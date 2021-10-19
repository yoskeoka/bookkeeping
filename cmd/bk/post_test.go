package main

import "testing"

func Test_parseJournalItem(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantAccID  int
		wantAmount int
		wantDesc   string
		wantErr    bool
	}{
		{"ok, without description", args{"23/5000"},
			23, 5000, "", false,
		},
		{"ok, with description", args{"23/5000/foo bar"},
			23, 5000, "foo bar", false,
		},
		{"error, missing amount and separator", args{"23"},
			0, 0, "", true,
		},
		{"error, missing amount", args{"23/"},
			0, 0, "", true,
		},
		{"error, missing id", args{"/9999"},
			0, 0, "", true,
		},
		{"error, wrong id format", args{"abc/9999"},
			0, 0, "", true,
		},
		{"error, wrong amount format", args{"12/9999ab"},
			0, 0, "", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccID, gotAmount, gotDesc, err := parseJournalItem(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJournalItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAccID != tt.wantAccID {
				t.Errorf("parseJournalItem() gotAccID = %v, want %v", gotAccID, tt.wantAccID)
			}
			if gotAmount != tt.wantAmount {
				t.Errorf("parseJournalItem() gotAmount = %v, want %v", gotAmount, tt.wantAmount)
			}
			if gotDesc != tt.wantDesc {
				t.Errorf("parseJournalItem() gotDesc = %v, want %v", gotDesc, tt.wantDesc)
			}
		})
	}
}
