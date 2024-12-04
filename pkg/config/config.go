package config

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
)

type Context struct {
    Name  string `json:"name"`
    URL   string `json:"url"`
    Token string `json:"token"`
}

type Config struct {
    CurrentContext string    `json:"current_context"`
    Contexts       []Context `json:"contexts"`
}

var configFilePath string

func init() {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error getting home directory:", err)
        os.Exit(1)
    }
    configDir := filepath.Join(homeDir, ".dtctl")
    if _, err := os.Stat(configDir); os.IsNotExist(err) {
        os.Mkdir(configDir, 0755)
    }
    configFilePath = filepath.Join(configDir, "config.json")
}

func loadConfig() (*Config, error) {
    var cfg Config
    if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
        return &cfg, nil
    }
    data, err := ioutil.ReadFile(configFilePath)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &cfg)
    if err != nil {
        return nil, err
    }
    return &cfg, nil
}

func saveConfig(cfg *Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(configFilePath, data, 0644)
    if err != nil {
        return err
    }
    return nil
}

func AddContext(ctx Context) error {
    cfg, err := loadConfig()
    if err != nil {
        return err
    }
    for _, c := range cfg.Contexts {
        if c.Name == ctx.Name {
            return fmt.Errorf("context '%s' already exists", ctx.Name)
        }
    }
    cfg.Contexts = append(cfg.Contexts, ctx)
    return saveConfig(cfg)
}

func UseContext(name string) error {
    cfg, err := loadConfig()
    if err != nil {
        return err
    }
    found := false
    for _, c := range cfg.Contexts {
        if c.Name == name {
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("context '%s' not found", name)
    }
    cfg.CurrentContext = name
    return saveConfig(cfg)
}

func GetConfig() (*Config, error) {
    return loadConfig()
}

func GetContext(name string) (*Context, error) {
    cfg, err := loadConfig()
    if err != nil {
        return nil, err
    }
    for _, ctx := range cfg.Contexts {
        if ctx.Name == name {
            return &ctx, nil
        }
    }
    return nil, fmt.Errorf("context '%s' not found", name)
}

func GetCurrentContext() (*Context, error) {
    cfg, err := loadConfig()
    if err != nil {
        return nil, err
    }
    if cfg.CurrentContext == "" {
        return nil, fmt.Errorf("no current context is set")
    }
    return GetContext(cfg.CurrentContext)
}

func UpdateContext(updatedCtx Context) error {
    cfg, err := loadConfig()
    if err != nil {
        return err
    }
    found := false
    for i, ctx := range cfg.Contexts {
        if ctx.Name == updatedCtx.Name {
            cfg.Contexts[i] = updatedCtx
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("context '%s' not found", updatedCtx.Name)
    }
    return saveConfig(cfg)
}
