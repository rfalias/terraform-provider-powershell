package pypwsh

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/rfalias/gopypwsh"

	"os"
	"time"
	"math/rand"
)

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func waitForLock(client *Powershell) bool {
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(100)) * time.Millisecond)

	locked := fileExists(client.lockfile)

	for locked == true {
		time.Sleep(100 * time.Millisecond)
		locked = fileExists(client.lockfile)
	}

	time.Sleep(1000 * time.Millisecond)
	return true
}

func resourcePyPwsh() *schema.Resource {
	return &schema.Resource{
		Create: resourcePyPwshRecordCreate,
		Read:   resourcePyPwshRecordRead,
		Delete: resourcePyPwshRecordDelete,

		Schema: map[string]*schema.Schema{
			"cmd": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePyPwshRecordCreate(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*Powershell)

	cmd := d.Get("cmd").(string)


	waitForLock(client)
	
	file, err := os.Create(client.lockfile)
	if err != nil {
		return err
	}

        var id string = cmd
	_, err = gopypwsh.RunPyCommandCreate(client.username, client.password, client.server, cmd, client.py)

	if err != nil {
		//something bad happened
		return err
	}

	d.SetId(id)

	file.Close()
	os.Remove(client.lockfile)

	return nil
}


func resourcePyPwshRecordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePyPwshRecordDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}
