package command

import (
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/spf13/cobra"

	"github.com/juancwu/konbini/shared/form"
)

var membershipCmd = &cobra.Command{
	Use:   "membership",
	Long:  "Become a member of Konbini to gain access to all the awesome services.",
	Short: "Become a member of Konbini.",
	Run:   membershipRun,
}

var membershipForm *form.MembershipForm

func init() {
	membershipForm = new(form.MembershipForm)
	membershipCmd.PersistentFlags().StringVar(&membershipForm.Email, "email", "", "Email to link to membership")
	membershipCmd.PersistentFlags().StringVar(&membershipForm.Password, "password", "", "Passowrd for membership")
	membershipCmd.PersistentFlags().StringVar(&membershipForm.FirstName, "firstname", "", "Your first name")
	membershipCmd.PersistentFlags().StringVar(&membershipForm.LastName, "lastname", "", "Your last name")
}

func membershipRun(cmd *cobra.Command, args []string) {
	if !prompt {
		log.Debug("No prompt. Getting values from flags", "cmd", "konbini get membership")
		validate := validator.New()
		if err := validate.Struct(membershipForm); err != nil {
			log.Fatalf("One or more fields are invalid/missing: %v\n", err)
		}
	} else {
		log.Debug("Building prompt...", "cmd", "konbini get membership")
	}
}
