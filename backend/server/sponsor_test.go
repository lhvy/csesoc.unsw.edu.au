package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

const companyName = "Example"
const companyLogo = "https://static.canva.com/static/images/canva_logo_100x100@2x.png"
const companyTier = "2"
const companyDetail = "Example"
const getRequest = "http://localhost:1323/api/v1/sponsors/"

func TestSponsor(t *testing.T) {

	t.Run("Sponsor setup test", func(t *testing.T) {
		resp, err := http.Get("http://localhost:1323/api/v1/sponsors")
		if err != nil {
			t.Errorf("Could not get perform request.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusOK)
		var sponsors []*Sponsor
		if err = json.NewDecoder(resp.Body).Decode(&sponsors); err != nil {
			t.Errorf("Error parsing json response: %s", err)
		}
		if len(sponsors) == 0 {
			t.Errorf("Sponsors were not populated.")
		}
	})

	t.Run("Testing sponsor filtering", func(t *testing.T) {
		resp, err := http.Get("http://localhost:1323/api/v1/sponsors?tier=2")
		if err != nil {
			t.Errorf("Could not perform get sponsors request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusOK)
	})

	t.Run("New sponsor", func(t *testing.T) {
		client := &http.Client{}
		form := url.Values{
			"name":   {companyName},
			"logo":   {companyLogo},
			"tier":   {companyTier},
			"detail": {companyDetail},
		}
		req, _ := http.NewRequest("POST", "http://localhost:1323/api/v1/sponsors", strings.NewReader(form.Encode()))
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTkzMTI5NjAwLCJ6SUQiOiJ6NTEyMzQ1NiJ9.jYC2qlpzAKIMPFywQ6pWIV1qat_h7OrorJ-zQM5jDpg")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		// resp, err := http.PostForm("http://localhost:1323/api/v1/sponsors", url.Values{
		// 	"name":   {companyName},
		// 	"logo":   {companyLogo},
		// 	"tier":   {companyTier},
		// 	"detail": {companyDetail},
		// })
		if err != nil {
			t.Errorf("Could not perform post sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusCreated)
	})

	t.Run("Get newly created sponsor", func(t *testing.T) {
		resp, err := http.Get(getRequest + companyName)
		if err != nil {
			t.Errorf("Could not perform get sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusOK)

		var newSponsor *Sponsor
		if err = json.NewDecoder(resp.Body).Decode(&newSponsor); err != nil {
			t.Errorf("Error parsing json response: %s", err)
		} else {
			assertResponseBody(t, newSponsor.Name, companyName)
			assertResponseBody(t, newSponsor.Logo, companyLogo)
			assertResponseBody(t, strconv.Itoa(newSponsor.Tier), companyTier)
		}
	})

	t.Run("Delete newly created sponsor", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", getRequest+companyName, nil)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTkzMTI5NjAwLCJ6SUQiOiJ6NTEyMzQ1NiJ9.jYC2qlpzAKIMPFywQ6pWIV1qat_h7OrorJ-zQM5jDpg")
		if err != nil {
			t.Errorf("Could not create delete request for sponsor.")
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Could not perform delete sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusNoContent)
	})

	t.Run("Check newly removed sponsor", func(t *testing.T) {
		resp, err := http.Get(getRequest + companyName)
		if err != nil {
			t.Errorf("Could not perform get sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusNotFound)
	})
}

func TestSponsorError(t *testing.T) {
	t.Run("Duplicate Sponsor sponsor", func(t *testing.T) {
		client := &http.Client{}
		form := url.Values{
			"name":   {companyName},
			"logo":   {companyLogo},
			"tier":   {companyTier},
			"detail": {companyDetail},
		}
		req, _ := http.NewRequest("POST", "http://localhost:1323/api/v1/sponsors", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTkzMTI5NjAwLCJ6SUQiOiJ6NTEyMzQ1NiJ9.jYC2qlpzAKIMPFywQ6pWIV1qat_h7OrorJ-zQM5jDpg")
		req.PostForm = form
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Could not perform post sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusCreated)

		req, _ = http.NewRequest("POST", "http://localhost:1323/api/v1/sponsors", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTkzMTI5NjAwLCJ6SUQiOiJ6NTEyMzQ1NiJ9.jYC2qlpzAKIMPFywQ6pWIV1qat_h7OrorJ-zQM5jDpg")
		req.PostForm = form
		resp, err = client.Do(req)
		if err != nil {
			t.Errorf("Could not perform post sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusConflict)

		req, err = http.NewRequest("DELETE", getRequest+companyName, nil)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTkzMTI5NjAwLCJ6SUQiOiJ6NTEyMzQ1NiJ9.jYC2qlpzAKIMPFywQ6pWIV1qat_h7OrorJ-zQM5jDpg")
		if err != nil {
			t.Errorf("Could not create delete request for sponsor.")
		}
		resp, err = client.Do(req)
		if err != nil {
			t.Errorf("Could not perform delete sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusNoContent)
	})

	t.Run("Missing parameters when creating", func(t *testing.T) {
		client := &http.Client{}
		req, _ := http.NewRequest("POST", "http://localhost:1323/api/v1/sponsors", nil)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTkzMTI5NjAwLCJ6SUQiOiJ6NTEyMzQ1NiJ9.jYC2qlpzAKIMPFywQ6pWIV1qat_h7OrorJ-zQM5jDpg")
		form := url.Values{
			"name": {companyName},
			"logo": {companyLogo},
		}
		req.PostForm = form
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Could not perform post sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusBadRequest)
	})

	t.Run("Get non existent sponsor", func(t *testing.T) {
		resp, err := http.Get(getRequest + "nonexistent")
		if err != nil {
			t.Errorf("Could not perform get sponsor request. Check connection.")
		}
		defer resp.Body.Close()

		assertStatus(t, resp.StatusCode, http.StatusNotFound)
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Response body is wrong, got %s, want %s", got, want)
	}
}