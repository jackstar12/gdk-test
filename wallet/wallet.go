package wallet

/*
#cgo CFLAGS: -I${SRCDIR}/include/
#cgo LDFLAGS: -L${SRCDIR}/lib -lgreenaddress -Wl,-rpath=${SRCDIR}/lib
#include "gdk.h"
#include "stdio.h"
#include "stdlib.h"
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

type AuthHandler = *C.struct_GA_auth_handler
type Json = *C.GA_json

type Wallet struct {
	session *C.struct_GA_session
}

func toErr(ret C.int, errMessage string) error {
	if ret == C.GA_OK {
		return nil
	}
	return fmt.Errorf("failed with code %v: %v", ret, errMessage)
}

func toJson(data interface{}) (result Json, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	string := C.CString(string(bytes))
	defer C.free(unsafe.Pointer(string))
	err = toErr(C.GA_convert_string_to_json(string, &result), "failed to convert to json")
	return result, err
}

func withOutput(ret C.int, errMessage string, output Json) (interface{}, error) {
	err := toErr(ret, errMessage)
	if err != nil {
		return nil, err
	}
	string := C.CString("")
	defer C.free(unsafe.Pointer(string))
	defer C.GA_destroy_json(output)
	if err := toErr(C.GA_convert_json_to_string(output, &string), errMessage+": failed to parse json"); err != nil {
		return nil, err
	}

	var v interface{}

	fmt.Println(C.GoString(string))

	return v, json.Unmarshal([]byte(C.GoString(string)), &v)
}

func withAuthhandler(ret C.int, errMessage string, handler AuthHandler) (interface{}, error) {
	defer C.GA_destroy_auth_handler(handler)
	var output Json

	if err := toErr(ret, errMessage); err != nil {
		return nil, err
	}
	return withOutput(C.GA_auth_handler_get_status(handler, &output), errMessage+": auth handler status", output)
}

func (wallet *Wallet) Init() error {
	params, err := toJson(map[string]interface{}{
		"datadir":   "./data",
		"log_level": "debug",
	})
	defer C.GA_destroy_json(params)
	if err != nil {
		return err
	}
	if err := toErr(C.GA_init(params), "failed to initialize"); err != nil {
		return err
	}

	if err := toErr(C.GA_create_session(&wallet.session), "failed to create session"); err != nil {
		return err
	}

	params, err = toJson(map[string]interface{}{
		"name":      "testnet-liquid",
		"log_level": "debug",
	})

	if err := toErr(C.GA_connect(wallet.session, params), "failed to connect"); err != nil {
		return err
	}

	return err
}

func (wallet *Wallet) Register() (string, error) {
	buffer := C.CString("")
	defer C.free(unsafe.Pointer(buffer))
	if err := toErr(C.GA_generate_mnemonic(&buffer), "failed to generate mnemonic"); err != nil {
		return "", err
	}
	mnemonic := C.GoString(buffer)

	login, err := toJson(map[string]string{"mnemonic": mnemonic})
	if err != nil {
		return "", err
	}
	hwDevice, err := toJson(map[string]interface{}{})
	if err != nil {
		return "", err
	}

	var handler AuthHandler
	_, err = withAuthhandler(C.GA_register_user(wallet.session, hwDevice, login, &handler), "failed to register user", handler)
	if err != nil {
		return "", err
	}

	if err := wallet.Login(mnemonic); err != nil {
		return "", errors.New("successfully registered but failed to login: " + err.Error())
	}

	return mnemonic, nil
}

func (wallet *Wallet) Login(mnemonic string) error {
	fmt.Println("Logging with mnemonic: " + mnemonic)
	login, err := toJson(map[string]string{"mnemonic": mnemonic})
	if err != nil {
		return err
	}
	hwDevice, err := toJson(map[string]interface{}{})
	if err != nil {
		return err
	}

	var handler AuthHandler
	res, err := withAuthhandler(C.GA_login_user(wallet.session, hwDevice, login, &handler), "failed to register user", handler)
	if err != nil {
		return err
	}
	fmt.Println("logged in: ", res)
	return nil
}
