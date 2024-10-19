package entities

type InvitationToExpose struct {
    Uuid     string `json:"uuid"`
    Email    string `json:"email"`
    CreatedAt string `json:"created_at"`
    InvitationCode string `json:"invitation_code"`
}
