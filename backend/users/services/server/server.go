package server

import (
	"context"
	"net"
	"net/http"
	"ylem_users/api"
	"ylem_users/config"
	"ylem_users/helpers"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Listen string
}

func (s *Server) Run(ctx context.Context) error {
	log.Info("Starting server listening on " + s.Listen)

	rdbClient := helpers.RedisDbConn(ctx)
	dbSource := api.NewDbSource(rdbClient)
	authMiddleware := &api.AuthMiddleware{DbSource: dbSource, Ctx: ctx}

	rtr := mux.NewRouter()

	rtr.Handle("/organization/{uuid}/users", authMiddleware.AuthCheck(http.HandlerFunc(api.GetUsersInOrganization))).Methods(http.MethodGet)
	rtr.Handle("/organization/{uuid}/pending-invitations", authMiddleware.AuthCheck(http.HandlerFunc(api.GetPendingInvitationsInOrganization))).Methods(http.MethodGet)
	rtr.Handle("/organization/{uuid}/invitations", authMiddleware.AuthCheck(http.HandlerFunc(api.PostInvitations))).Methods(http.MethodPost)
	rtr.Handle("/invitations/{key}/validate", authMiddleware.IsNotAuthCheck(http.HandlerFunc(api.ValidateInvitation))).Methods(http.MethodPost)

	rtr.Handle("/user/{uuid}/delete", authMiddleware.AuthCheck(http.HandlerFunc(api.DeleteUser))).Methods(http.MethodPost)
	rtr.Handle("/user/{uuid}/assign-role", authMiddleware.AuthCheck(http.HandlerFunc(api.AssignRoleToUser))).Methods(http.MethodPost)
	rtr.Handle("/user/{uuid}/activate", authMiddleware.AuthCheck(http.HandlerFunc(api.ActivateUser))).Methods(http.MethodPost)
	rtr.Handle("/user", authMiddleware.IsNotAuthCheck(http.HandlerFunc(authMiddleware.PostUser))).Methods(http.MethodPost)
	rtr.Handle("/login", authMiddleware.IsNotAuthCheck(http.HandlerFunc(authMiddleware.LoginUser))).Methods(http.MethodPost)
	rtr.Handle("/logout", authMiddleware.AuthCheck(http.HandlerFunc(authMiddleware.LogoutUser))).Methods(http.MethodPost)
	rtr.Handle("/me/password", authMiddleware.AuthCheck(http.HandlerFunc(api.UpdateMyPassword))).Methods(http.MethodPost)
	rtr.Handle("/me", authMiddleware.AuthCheck(http.HandlerFunc(api.UpdateMe))).Methods(http.MethodPost)
	rtr.Handle("/organization/{uuid}", authMiddleware.AuthCheck(http.HandlerFunc(api.UpdateOrganization))).Methods(http.MethodPost)
	rtr.Handle("/my-organization", authMiddleware.AuthCheck(http.HandlerFunc(api.GetMyOrganization))).Methods(http.MethodGet)
	rtr.Handle("/email/{key}/confirm", http.HandlerFunc(api.ConfirmEmail)).Methods(http.MethodPost)

	rtr.Handle("/auth/{provider}", http.HandlerFunc(api.ExternalAuth)).Methods(http.MethodGet)
	rtr.Handle("/auth/{provider}/available", http.HandlerFunc(api.IsExternalAuthAvailable)).Methods(http.MethodGet)
	rtr.Handle("/auth/{provider}/callback", authMiddleware.IsNotAuthCheck(http.HandlerFunc(authMiddleware.ExternalAuthCallback))).Methods(http.MethodPost)

	// The following endpoints are for the internal use only and should be
	// protected on the server level from any outside of the network use
	rtr.HandleFunc("/private/user/check-permission", api.PermissionCheck).Methods(http.MethodPost)
	rtr.Handle("/private/user/check-authorization", authMiddleware.AuthCheck(http.HandlerFunc(api.AuthorizationCheck))).Methods(http.MethodPost)
	rtr.HandleFunc("/private/organization/{uuid}/update-connections", api.UpdateOrganizationConnections).Methods(http.MethodPost)
	rtr.Handle("/private/user/{uuid}/confirm-email", http.HandlerFunc(api.ConfirmEmailInternal)).Methods(http.MethodPost)
	rtr.Handle("/private/organization/{uuid}/data-key", http.HandlerFunc(api.GetOrganizationDataKey)).Methods(http.MethodGet)

	rtr.HandleFunc("/private/jwt-tokens/", authMiddleware.IssueJWTPrivate).Methods(http.MethodPost)

	http.Handle("/", rtr)

	server := &http.Server{
		Addr:    config.Cfg().Listen,
		Handler: nil,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(ctx)
		rdbClient.Close()
	}()

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func NewServer(listen string) *Server {
	s := &Server{
		Listen: listen,
	}

	return s
}
