package csv

import "testing"

func TestValidateHeader(t *testing.T) {
	t.Run("valid header", func(t *testing.T) {
		content := []byte("role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,1000\n")
		if err := ValidateHeader(content); err != nil {
			t.Fatalf("expected valid header, got error: %v", err)
		}
	})

	t.Run("invalid header", func(t *testing.T) {
		content := []byte("role,company,location,link,salary\nGo Dev,Acme,BR,https://acme.com,1000\n")
		if err := ValidateHeader(content); err == nil {
			t.Fatalf("expected invalid header error")
		}
	})
}

func TestParseAndValidate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		content := []byte("role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,1000\n")
		parsed, rowErrors, err := ParseAndValidate(content)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(rowErrors) != 0 {
			t.Fatalf("expected no row errors, got %d", len(rowErrors))
		}
		if len(parsed) != 1 {
			t.Fatalf("expected 1 parsed row, got %d", len(parsed))
		}
		if parsed[0].LineNumber != 2 {
			t.Fatalf("expected line number 2, got %d", parsed[0].LineNumber)
		}
	})

	t.Run("invalid row", func(t *testing.T) {
		content := []byte("role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,0\n")
		parsed, rowErrors, err := ParseAndValidate(content)
		if err != nil {
			t.Fatalf("expected no parse error, got %v", err)
		}
		if len(parsed) != 0 {
			t.Fatalf("expected no parsed rows, got %d", len(parsed))
		}
		if len(rowErrors) != 1 {
			t.Fatalf("expected 1 row error, got %d", len(rowErrors))
		}
		if rowErrors[0].LineNumber != 2 {
			t.Fatalf("expected error in line 2, got %d", rowErrors[0].LineNumber)
		}
	})
}
