package pypwsh

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"fmt"
	"os"
	"io/ioutil"
)

// Provider allows making changes to Windows DNS server
// Utilises Powershell to connect to domain controller
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("USERNAME", nil),
				Description: "Username to connect to AD.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PASSWORD", nil),
				Description: "The password to connect to AD.",
			},
			"server": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVER", nil),
				Description: "The AD server to connect to.",
			},
			"cmd": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CMD", false),
				Description: "Powershell Command",
			},
                        "py": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("PY", false),
                                Description: "Path to the python powershell wrapper file",
                        },
		},
		ResourcesMap: map[string]*schema.Resource{
			"pypwsh": resourcePyPwsh(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, fmt.Errorf("The 'username' property was not specified.")
	}
	
        password := d.Get("password").(string)
        if password == ""{
                return nil, fmt.Errorf("The 'password' property was not specified.")
        }

	server := d.Get("server").(string)
	if server == "" {
		return nil, fmt.Errorf("The 'server' property was not specified.")
	}

	cmd := d.Get("cmd").(string)
	py := d.Get("py").(string)

	f, err := ioutil.TempFile("", "terraform-pypwsh")
	lockfile := f.Name()
	os.Remove(f.Name())

	client := Powershell {
		username:	username,
		password:	password,
		server:		server,
                cmd:            cmd,
                py:             py,
                lockfile:       lockfile,
	}

	return &client, err
}

type Powershell struct {
	username	string
	password	string
	server		string
        cmd		string
        py              string
        lockfile        string
}
