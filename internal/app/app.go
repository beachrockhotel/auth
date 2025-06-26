package app

import (
	"context"
	descAccess "github.com/beachrockhotel/auth/pkg/access_v1"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/beachrockhotel/auth/internal/closer"
	"github.com/beachrockhotel/auth/internal/config"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
	_ "github.com/beachrockhotel/auth/statik"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"

	"crypto/tls"
	"google.golang.org/grpc/credentials"

	"go.uber.org/zap/zapcore"

	"github.com/beachrockhotel/auth/internal/app"
	"github.com/beachrockhotel/auth/internal/logger"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

// NewApp создаёт и инициализирует все зависимости
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает GRPC и HTTP серверы параллельно
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := a.runGRPCServer(); err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runHTTPServer(); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Wait()
	return nil
}

// initDeps инициализирует все зависимости по порядку
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	return config.Load(".env")
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	cert, err := tls.LoadX509KeyPair("service.pem", "service.key")
	if err != nil {
		return err
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})

	a.grpcServer = grpc.NewServer(grpc.Creds(creds))
	reflection.Register(a.grpcServer)
	desc.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true, // на dev можно так, на prod надо валидный cert
	})

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds), // ✅ используем TLS
	}

	err := desc.RegisterAuthV1HandlerFromEndpoint(
		ctx,
		mux,
		a.serviceProvider.GRPCConfig().Address(),
		opts,
	)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile(statikFs, "/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	addr := a.serviceProvider.GRPCConfig().Address()
	logger.Info("GRPC server is running", zap.String("addr", addr))

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(listener)
}

func (a *App) runHTTPServer() error {
	addr := a.serviceProvider.HTTPConfig().Address()
	logger.Info("HTTP server is running", zap.String("addr", addr))

	return a.httpServer.ListenAndServe()

}

func (a *App) runSwaggerServer() error {
	addr := a.serviceProvider.SwaggerConfig().Address()
	logger.Info("Swagger server is running", zap.String("addr", addr))

	return a.swaggerServer.ListenAndServe()
}

func serveSwaggerFile(fs http.FileSystem, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Serving swagger file", zap.String("path", path))

		file, err := fs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Error opening swagger file", zap.Error(err))
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Error reading swagger file", zap.Error(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Error writing swagger response", zap.Error(err))
			return
		}

		logger.Info("Successfully served swagger file", zap.String("path", path))
	}
}
