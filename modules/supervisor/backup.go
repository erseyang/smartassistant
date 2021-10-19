package supervisor

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/zhiting-tech/smartassistant/modules/config"

	"github.com/sirupsen/logrus"

	"github.com/zhiting-tech/smartassistant/modules/plugin/docker"

	"github.com/zhiting-tech/smartassistant/modules/plugin"

	jsoniter "github.com/json-iterator/go"

	"github.com/zhiting-tech/smartassistant/modules/types"
)

// Backup 备份描述文件结构 backup.json
type Backup struct {
	Name           string         `json:"name"`
	Date           time.Time      `json:"date"`
	SmartAssistant SmartAssistant `json:"smartassistant"`
	Plugins        []Plugin       `json:"plugins"`
}

type SmartAssistant struct {
	Version string `json:"version"`
}

func (s SmartAssistant) RefStr() string {
	img := docker.Image{
		Name:     saImage.Name,
		Tag:      s.Version,
		Registry: saImage.Registry,
	}
	return img.RefStr()
}

type Plugin struct {
	ID      string       `json:"id"`
	Brand   string       `json:"brand"`
	Image   docker.Image `json:"image"`
	Version string       `json:"version"`
}

func newBackup(name string) *Backup {
	var plugins []Plugin
	if plgs, err := plugin.GetGlobalManager().Load(); err == nil {
		plugins = make([]Plugin, 0, len(plgs))
		for _, plg := range plgs {
			plugins = append(plugins, Plugin{
				ID:      plg.ID,
				Brand:   plg.Brand,
				Image:   plg.Image,
				Version: plg.Version,
			})
		}
	}
	return &Backup{
		Name: name,
		Date: time.Now(),
		SmartAssistant: SmartAssistant{
			Version: types.Version,
		},
		Plugins: plugins,
	}
}

func loadBackup(dir string) *Backup {
	var b Backup
	d, err := ioutil.ReadFile(filepath.Join(dir, "backup.json"))
	if err != nil {
		return nil
	}
	err = jsoniter.Unmarshal(d, &b)
	if err != nil {
		return nil
	}
	return &b
}

func (b *Backup) Save(dir string) (err error) {
	if err = b.writeBackupJson(dir); err != nil {
		return
	}
	if err = b.backupImages(dir); err != nil {
		return
	}
	if err = b.backupRuntimeDir(dir); err != nil {
		return
	}
	return
}

func (b *Backup) writeBackupJson(dir string) (err error) {
	c, err := jsoniter.MarshalIndent(*b, "", "  ")
	if err != nil {
		return
	}
	f, err := os.Create(filepath.Join(dir, "backup.json"))
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(c)
	return
}

// 导出 smart assistant 以及插件的镜像
func (b *Backup) backupImages(target string) (err error) {
	cli := docker.GetClient()
	plgs, err := plugin.GetGlobalManager().Load()
	if err != nil {
		return
	}
	imgs := make([]docker.Image, 0, len(plgs)+1)
	imgs = append(imgs, saImage)
	for _, plg := range plgs {
		imgs = append(imgs, plg.Image)
	}
	for _, img := range imgs {
		if !cli.IsImageAdd(img.RefStr()) {
			logrus.Infof("image %v not added, pulling", img.RefStr())
			if err = cli.Pull(img.RefStr()); err != nil {
				return
			}
		}
	}
	err = cli.ImageSave(filepath.Join(target, "images.tar"), imgs...)
	return
}

// 导出数据,配置目录
func (b *Backup) backupRuntimeDir(target string) (err error) {
	// copy dir
	runtimeDir := config.GetConf().SmartAssistant.RuntimePath
	logrus.Infof("backup %v", runtimeDir)

	dcSrcFile := filepath.Join(runtimeDir, "docker-compose.yaml")
	dcSrcInfo, err := os.Stat(dcSrcFile)
	if err != nil {
		return
	}
	dcDstFile := filepath.Join(target, "docker-compose.yaml")
	err = copyFile(dcSrcFile, dcDstFile, dcSrcInfo)
	if err != nil {
		return
	}

	dataSrcDir := filepath.Join(runtimeDir, "data")
	dataSrcInfo, err := os.Stat(dataSrcDir)
	if err != nil {
		return
	}
	dataDstDir := filepath.Join(target, "data")
	os.MkdirAll(dataDstDir, dataSrcInfo.Mode())
	if err = copyDir(dataSrcDir, dataDstDir); err != nil {
		return
	}

	configSrcDir := filepath.Join(runtimeDir, "config")
	configSrcInfo, err := os.Stat(configSrcDir)
	if err != nil {
		return
	}
	configDstDir := filepath.Join(target, "config")
	os.MkdirAll(configDstDir, configSrcInfo.Mode())
	if err = copyDir(configSrcDir, configDstDir); err != nil {
		return
	}

	return
}

func copyDir(srcDir string, dstDir string) (err error) {
	err = filepath.Walk(srcDir, func(srcPath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)
		switch mode := fi.Mode(); {
		case mode.IsDir():
			if err := os.Mkdir(dstPath, fi.Mode()); err != nil && !os.IsExist(err) {
				return err
			}
		case mode.IsRegular():
			if err := copyFile(srcPath, dstPath, fi); err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func copyFile(srcPath string, dstPath string, fi os.FileInfo) (err error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fi.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return

}
