package gateway
 
import (
	"fmt"
	"net/http"
	"path"
	"strings"
 
	assetfs "github.com/elazarl/go-bindata-assetfs"
 
	swagger "hello/gateway/swagger"
	gw "hello/protos"
 
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)
 
///
func HttpRun(gprcPort, httpPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
 
	gwmux, err := newGateway(ctx, gprcPort)
	if err != nil {
		panic(err)
	}
 
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.HandleFunc("/swagger/", serveSwaggerFile)
	serveSwaggerUI(mux)
 
	fmt.Println("grpc-gateway listen on localhost" + httpPort)
	return http.ListenAndServe(httpPort, mux)
}
 
func newGateway(ctx context.Context, gprcPort string) (http.Handler, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}
 
	gwmux := runtime.NewServeMux()
	if err := gw.RegisterGreeterHandlerFromEndpoint(ctx, gwmux, gprcPort, opts); err != nil {
		return nil, err
	}
 
	return gwmux, nil
}
 
func serveSwaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		fmt.Printf("Not Found: %s\r\n", r.URL.Path)
		http.NotFound(w, r)
		return
	}
 
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join("../protos", p)
 
	fmt.Printf("Serving swagger-file: %s\r\n", p)
 
	http.ServeFile(w, r, p)
}
 
func serveSwaggerUI(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "../third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}