package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Response struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Model   string `json:"model"`
	Name    string `json:"name"`
	Key     string `json:"key"`
	Account struct {
		ID                   string `json:"id"`
		PrivateKey           string `json:"private_key"`
		ReservedHex          string `json:"reserved_hex"`
		ReservedDec          []int  `json:"reserved_dec"`
		AccountType          string `json:"account_type"`
		Created              string `json:"created"`
		Updated              string `json:"updated"`
		PremiumData          int    `json:"premium_data"`
		Quota                int    `json:"quota"`
		Usage                int    `json:"usage"`
		WarpPlus             bool   `json:"warp_plus"`
		ReferralCount        int    `json:"referral_count"`
		ReferralRenewalCount int    `json:"referral_renewal_countdown"`
		Role                 string `json:"role"`
		License              string `json:"license"`
	} `json:"account"`
	Config struct {
		ClientID string `json:"client_id"`
		Peers    []struct {
			PublicKey string `json:"public_key"`
			Endpoint  struct {
				V4   string `json:"v4"`
				V6   string `json:"v6"`
				Host string `json:"host"`
			} `json:"endpoint"`
		} `json:"peers"`
		Interface struct {
			Addresses struct {
				V4 string `json:"v4"`
				V6 string `json:"v6"`
			} `json:"addresses"`
		} `json:"interface"`
		Services struct {
			HTTPProxy string `json:"http_proxy"`
		} `json:"services"`
	} `json:"config"`
	Token     string `json:"token"`
	Warp      bool   `json:"warp_enabled"`
	Waitlist  bool   `json:"waitlist_enabled"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	TOS       string `json:"tos"`
	Place     int    `json:"place"`
	Locale    string `json:"locale"`
	Enabled   bool   `json:"enabled"`
	InstallID string `json:"install_id"`
	FCMToken  string `json:"fcm_token"`
	SerialNum string `json:"serial_number"`
}

func main() {
	action := ParseCommandLine()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if action.Help {
		help()
		return
	}

	if action.Register {
		store, output, err := register()
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
		if action.FileName != "" {
			err = os.WriteFile(action.FileName, store, 0600)
			if err != nil {
				panic(err)
			}
			return
		} else {
			fileName := "wgcf.json"
			editedFileName := "wgcf.json"
			i := 0

			for {
				if _, err := os.Stat(fileName); err == nil {
					fileName = fmt.Sprintf("%s-%d.json", editedFileName[:len(editedFileName)-5], i)
					i++
				} else {
					break
				}
			}

			err := os.WriteFile(fileName, store, 0600)
			if err != nil {
				panic(err)
			}
			return
		}
	}

	if action.FileName == "" {
		action.FileName = "wgcf.json"
	} else if strings.HasPrefix(action.FileName, "-") {
		err := fmt.Sprintln("The parameter must not start with '-'")
		panic(err)
	}

	if !action.Bind && !action.UnBind && !action.Cancle && action.License == "" && action.Name == "" {
		err := fmt.Sprintln("You need to specify an action")
		panic(err)
	}

	if action.Bind {
		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := getBindingDevices(token, id)
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
		return
	}

	if action.UnBind {
		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := unBind(token, id)
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
		return
	}

	if action.Cancle {
		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		err = cancleAccount(token, id)
		if err != nil {
			panic(err)
		}

		err = os.Remove(action.FileName)
		if err != nil {
			panic(err)
		}
		fmt.Println("Cancled")
		return
	}

	if action.License != "" {

		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := changeLicense(token, id, action.License)
		if err != nil {
			panic(err)
		}
		fmt.Println(output)

		err = editFile(action.FileName, action.License)
		if err != nil {
			panic(err)
		}
		return
	}

	if action.Name != "" && !strings.HasPrefix(action.Name, "-") {
		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := changeName(token, id, action.Name)
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
		return
	} else {
		err := fmt.Sprintln("The parameter must not start with '-'")
		panic(err)
	}
}
