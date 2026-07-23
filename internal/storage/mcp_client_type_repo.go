package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MCPClientType struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ConfigKey       string `json:"config_key"`
	ConfigFormat    string `json:"config_format"`
	ConfigTemplate  string `json:"config_template"`
	FileName        string `json:"file_name"`
	FileDescription string `json:"file_description"`
	SortOrder       int    `json:"sort_order"`
	IsActive        bool   `json:"is_active"`
}

type MCPClientTypeRepository struct {
	db *pgxpool.Pool
}

func NewMCPClientTypeRepository(db *pgxpool.Pool) *MCPClientTypeRepository {
	return &MCPClientTypeRepository{db: db}
}

func (r *MCPClientTypeRepository) ListActive() ([]MCPClientType, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, name, config_key, config_format, config_template, file_name, file_description, sort_order, is_active
		 FROM mcp_client_types WHERE is_active = true ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MCPClientType
	for rows.Next() {
		var t MCPClientType
		if err := rows.Scan(&t.ID, &t.Name, &t.ConfigKey, &t.ConfigFormat, &t.ConfigTemplate,
			&t.FileName, &t.FileDescription, &t.SortOrder, &t.IsActive); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *MCPClientTypeRepository) Seed() error {
	type seed struct {
		name, configKey, configFormat, configTemplate, fileName, fileDescription string
		sortOrder                                                                int
	}

	seeds := []seed{
		{
			name: "Claude Desktop", configKey: "mcpServers", configFormat: "json", sortOrder: 1,
			fileName: "claude_desktop_config.json",
			fileDescription: "~/Library/Application Support/Claude/ (macOS) или %APPDATA%\\Claude\\ (Windows)",
			configTemplate: `{
  "mcpServers": {
    "mcp-is": {
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
		{
			name: "Cursor", configKey: "mcpServers", configFormat: "json", sortOrder: 2,
			fileName: "mcp.json",
			fileDescription: "~/.cursor/mcp.json",
			configTemplate: `{
  "mcpServers": {
    "mcp-is": {
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
		{
			name: "VS Code (Copilot)", configKey: "servers", configFormat: "json", sortOrder: 3,
			fileName: "mcp.json",
			fileDescription: ".vscode/mcp.json (в проекте)",
			configTemplate: `{
  "servers": {
    "mcp-is": {
      "type": "http",
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
		{
			name: "MiMoCode", configKey: "mcp", configFormat: "json", sortOrder: 4,
			fileName: "mcp.json",
			fileDescription: ".mimocode/mcp.json (в проекте)",
			configTemplate: `{
  "mcp": {
    "mcp-is": {
      "type": "remote",
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
		{
			name: "OpenCode", configKey: "mcpServers", configFormat: "json", sortOrder: 5,
			fileName: "config.json",
			fileDescription: ".opencode/config.json",
			configTemplate: `{
  "mcpServers": {
    "mcp-is": {
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
		{
			name: "LibreChat", configKey: "mcpServers", configFormat: "yaml", sortOrder: 6,
			fileName: "librechat.yaml",
			fileDescription: "корень проекта LibreChat",
			configTemplate: `mcpServers:
  mcp-is:
    type: streamable-http
    url: {{URL}}
    headers:
      Authorization: Bearer {{API_KEY}}`,
		},
		{
			name: "Open WebUI", configKey: "mcpServers", configFormat: "json", sortOrder: 7,
			fileName: "",
			fileDescription: "Через Admin Panel → Settings → MCP",
			configTemplate: `{
  "mcpServers": {
    "mcp-is": {
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
		{
			name: "Cline", configKey: "mcpServers", configFormat: "json", sortOrder: 8,
			fileName: "cline_mcp_settings.json",
			fileDescription: "%APPDATA%\\cline\\ (Windows) или ~/.config/cline/ (Linux/macOS)",
			configTemplate: `{
  "mcpServers": {
    "mcp-is": {
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
		{
			name: "Cherry Studio", configKey: "mcpServers", configFormat: "json", sortOrder: 9,
			fileName: "",
			fileDescription: "Через Settings → MCP Servers",
			configTemplate: `{
  "mcpServers": {
    "mcp-is": {
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
		{
			name: "Continue", configKey: "mcpServers", configFormat: "json", sortOrder: 10,
			fileName: "config.json",
			fileDescription: "~/.continue/config.json",
			configTemplate: `{
  "mcpServers": {
    "mcp-is": {
      "url": "{{URL}}",
      "headers": {
        "Authorization": "Bearer {{API_KEY}}"
      }
    }
  }
}`,
		},
	}

	for _, s := range seeds {
		_, err := r.db.Exec(context.Background(),
			`INSERT INTO mcp_client_types (name, config_key, config_format, config_template, file_name, file_description, sort_order)
			 SELECT $1, $2, $3, $4, $5, $6, $7
			 WHERE NOT EXISTS (SELECT 1 FROM mcp_client_types WHERE name = $1)`,
			s.name, s.configKey, s.configFormat, s.configTemplate, s.fileName, s.fileDescription, s.sortOrder,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
