package ai

import (
	"os"
	"testing"

)

func TestObtainDescription(t *testing.T) {


	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set")
	}

	imgURL := "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2dvLWltYWdlL2ltYWdlbmVzL21lMS5qcGc_WC1BbXotQWxnb3JpdGhtPUFXUzQtSE1BQy1TSEEyNTYmWC1BbXotQ3JlZGVudGlhbD1JOVYySk1FT1QwRUY3T1RLVFBRQyUyRjIwMjUxMDE3JTJGdXMtZWFzdC0xJTJGczMlMkZhd3M0X3JlcXVlc3QmWC1BbXotRGF0ZT0yMDI1MTAxN1QxODQ2MzRaJlgtQW16LUV4cGlyZXM9NDMyMDAmWC1BbXotU2VjdXJpdHktVG9rZW49ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmhZMk5sYzNOTFpYa2lPaUpKT1ZZeVNrMUZUMVF3UlVZM1QxUkxWRkJSUXlJc0ltVjRjQ0k2TVRjMk1EYzJPVEkwT0N3aWNHRnlaVzUwSWpvaWJXbHVhVzloWkcxcGJpSjkuQi1aekswMTI5MTZqZDdxa0p4NW44ckJqVmhoazVyazRkay1jc29NQzM3SmsyanVjTnhfaDJUbWp6MnlMdVVXbEtZQV9KLXJzbk4xUTlWc1U1RVBxREEmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JnZlcnNpb25JZD1udWxsJlgtQW16LVNpZ25hdHVyZT0xZTA0NTY0ODJkYzgxN2RmMGFiNDczYWEzNjcwMGZkMjNkNTEwY2NjM2I0NWRiYjEyOWM2ZWI5OWNmNGJjMjEy"
	desc, err := ObtainDescription(imgURL)
	if err != nil {
		t.Fatalf("Error calling ObtainDescription:", )
	}

	t.Logf(desc.Description, desc.Tags)

	t.Logf("Image description: %s", desc)
}
