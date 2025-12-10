package integration

import (
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errcode"
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errtext"
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errtype"
	"github.com/apotourlyan/ludus-studii/pkg/idutil"
	"github.com/apotourlyan/ludus-studii/pkg/testutil"
	"github.com/apotourlyan/ludus-studii/pkg/testutil/txutil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register"
	rerrcode "github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errcode"
	rerrtext "github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errtext"
)

func TestRegister_Success(t *testing.T) {
	t.Parallel()

	txID := idutil.UUID()

	txutil.TxTest(t, db, txID, func() {
		// Arrange
		path := getPath("/api/register")
		data := map[string]string{
			"email":    "test@example.com",
			"password": "Password123!",
		}

		// Act
		resp, code := testutil.SuccessRequest[register.Response](
			t, "POST", path, txID, data)

		// Assert response
		testutil.GotWant(t, code, 201)
		testutil.DontWantNil(t, resp)
		testutil.DontWant(t, resp.ID, 0)

		query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"
		exists := txutil.TxQueryValue[bool](t, txID, query, resp.ID)

		testutil.GotWant(t, exists, true)
	})
}

func TestRegister_EmailExists(t *testing.T) {
	t.Parallel()

	txID := idutil.UUID()

	txutil.TxTest(t, db, txID, func() {
		// Arrange
		path := getPath("/api/register")
		data := map[string]string{
			"email":    "test@example.com",
			"password": "Password123!",
		}

		// Act
		testutil.SuccessRequest[register.Response](
			t, "POST", path, txID, data)

		err, code := testutil.ErrorRequest[errorutil.ServiceErrorDto](
			t, "POST", path, txID, data)

		// Assert
		testutil.GotWant(t, code, 409)
		testutil.DontWantNil(t, err)
		testutil.GotWant(t, err.Type, errtype.Service)
		testutil.GotWant(t, err.Code, rerrcode.EmailExists)
		testutil.GotWant(t, err.Message, rerrtext.EmailExists)

		query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
		exists := txutil.TxQueryValue[bool](t, txID, query, "test@example.com")

		testutil.GotWant(t, exists, true)
	})
}

func TestRegister_InvalidRequest(t *testing.T) {
	t.Parallel()

	// Arrange
	path := getPath("/api/register")

	// Act
	err, code := testutil.ErrorRequest[errorutil.ServiceErrorDto](
		t, "POST", path, "", nil)

	testutil.GotWant(t, code, 400)
	testutil.DontWantNil(t, err)
	testutil.GotWant(t, err.Type, errtype.Service)
	testutil.GotWant(t, err.Code, errcode.Request)
	testutil.GotWant(t, err.Message, errtext.Request)
}

func TestRegister_DataRequired(t *testing.T) {
	t.Parallel()

	// Arrange
	path := getPath("/api/register")
	data := map[string]string{
		"email":    "",
		"password": "",
	}

	// Act
	err, code := testutil.ErrorRequest[errorutil.ValidationErrorDto](
		t, "POST", path, "", data)

	// Assert
	testutil.GotWant(t, code, 400)
	testutil.DontWantNil(t, err)
	testutil.GotWant(t, err.Type, errtype.Validation)

	got := err.Has("Email", errcode.Required, errtext.Required("Email"))
	testutil.GotWant(t, got, true)

	got = err.Has("Password", errcode.Required, errtext.Required("Password"))
	testutil.GotWant(t, got, true)
}
