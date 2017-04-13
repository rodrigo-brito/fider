package infra_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/infra"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

func TestIsAuthorized_WithAllowedRole(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetUser(&models.User{
		ID:   1,
		Role: models.RoleMember,
	})

	mw := infra.IsAuthorized(models.RoleAdministrator, models.RoleMember)
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestIsAuthorized_WithForbiddenRole(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetUser(&models.User{
		ID:   1,
		Role: models.RoleVisitor,
	})

	mw := infra.IsAuthorized(models.RoleAdministrator, models.RoleMember)
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusForbidden))
}

func TestIsAuthenticated_WithUser(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetUser(&models.User{
		ID: 1,
	})

	mw := infra.IsAuthenticated()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestIsAuthenticated_WithoutUser(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)

	mw := infra.IsAuthenticated()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusForbidden))
}

func TestHostChecker(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "login.test.canhearyou.com"

	mw := infra.HostChecker("http://login.test.canhearyou.com")
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestHostChecker_DifferentHost(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := infra.HostChecker("login.test.canhearyou.com")
	mw(app.HandlerFunc(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusBadRequest))
}
