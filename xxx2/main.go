package main

import (
	"fmt"
	"net/http"
)

const (
	err400 = /*0b00000001 */ 1 << (1 * iota)
	err401 /*= 0b00000010 */
	err404 /*= 0b00000100 */
	err406 /*= 0b00001000 */
	err422 /*= 0b00010000 */
	err500 /*= 0b00100000 */
	err501 /*= 0b01000000 */
	err502 /*= 0b10000000 */
)

func main() {

	fmt.Printf("err400 = % 4d %08b\n", err400, err400)
	fmt.Printf("err401 = % 4d %08b\n", err401, err401)
	fmt.Printf("err404 = % 4d %08b\n", err404, err404)
	fmt.Printf("err406 = % 4d %08b\n", err406, err406)
	fmt.Printf("err422 = % 4d %08b\n", err422, err422)
	fmt.Printf("err500 = % 4d %08b\n", err500, err500)
	fmt.Printf("err501 = % 4d %08b\n", err501, err501)
	fmt.Printf("err502 = % 4d %08b\n", err502, err502)

	codes := []int{500, 401, 300, 450, 422}

	x := throw(codes)
	fmt.Printf("%08b\n\n", x)

	if x&err400 == err400 {
		fmt.Println(http.StatusText(400))
	}
	if x&err401 == err401 {
		fmt.Println(http.StatusText(401))
	}
	if x&err404 == err404 {
		fmt.Println(http.StatusText(404))
	}
	if x&err406 == err406 {
		fmt.Println(http.StatusText(406))
	}
	if x&err422 == err422 {
		fmt.Println(http.StatusText(422))
	}
	if x&err500 == err500 {
		fmt.Println(http.StatusText(500))
	}
	if x&err501 == err501 {
		fmt.Println(http.StatusText(501))
	}
	if x&err502 == err502 {
		fmt.Println(http.StatusText(502))
	}

}

func throw(codes []int) int {
	var out int
	if contains(codes, http.StatusBadRequest) {
		out |= err400
	}
	if contains(codes, http.StatusUnauthorized) {
		out |= err401
	}
	if contains(codes, http.StatusNotFound) {
		out |= err404
	}
	if contains(codes, http.StatusNotAcceptable) {
		out |= err406
	}
	if contains(codes, http.StatusUnprocessableEntity) {
		out |= err422
	}
	if contains(codes, http.StatusInternalServerError) {
		out |= err500
	}
	if contains(codes, http.StatusBadGateway) {
		out |= err502
	}
	if contains(codes, http.StatusNotImplemented) {
		out |= err501
	}

	return out
}

func contains(s []int, x int) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}
