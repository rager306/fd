# M045 TEI Offline Startup Proof

Captured start: 2026-06-14T11:41:11Z

## Preflight

User explicitly confirmed controlled restart proof. Rollback if failure: remove `HF_HUB_OFFLINE=1` from `docker-compose.yaml`, run `docker compose up -d tei`, and re-run fd/TEI smoke.

Preflight container state:

```json
[
  {
    "name": "fd_tei",
    "status": "running",
    "health": "healthy",
    "started_at": "2026-06-14T09:24:48.743364832Z",
    "safe_env": [
      "HF_HOME=/data",
      "HUGGINGFACE_HUB_CACHE=/data",
      "PORT=80"
    ],
    "cmd": [
      "--model-id",
      "deepvk/USER-bge-m3"
    ]
  },
  {
    "name": "fd_api",
    "status": "running",
    "health": "healthy",
    "started_at": "2026-06-14T08:30:46.620802937Z",
    "safe_env": [
      "MODEL_ID=deepvk/USER-bge-m3",
      "PORT=8000",
      "TEI_URL=http://tei:80"
    ],
    "cmd": null
  },
  {
    "name": "fd_redis",
    "status": "running",
    "health": "healthy",
    "started_at": "2026-05-19T18:08:09.5269274Z",
    "safe_env": [],
    "cmd": [
      "redis-server",
      "--maxmemory",
      "2gb",
      "--maxmemory-policy",
      "allkeys-lru",
      "--save",
      "300",
      "1",
      "--appendonly",
      "no",
      "--protected-mode",
      "no"
    ]
  }
]
```

Compose candidate excerpt:

```text
name: fd
services:
  tei:
    command:
      - --model-id
      - deepvk/USER-bge-m3
    container_name: fd_tei
    environment:
      HF_HOME: /data
      HF_HUB_OFFLINE: "1"
      HUGGINGFACE_HUB_CACHE: /data
    healthcheck:
      test:
        - CMD
        - curl
        - -f
        - http://localhost:80/health
      timeout: 5s
      interval: 5s
      retries: 20
      start_period: 5m0s
    image: ghcr.io/huggingface/text-embeddings-inference:cpu-1.9
    networks:
      default: null
    ports:
      - mode: ingress
        target: 80
        published: "30080"
        protocol: tcp
    restart: unless-stopped
    volumes:
      - type: volume
        source: tei_data
        target: /data
        volume: {}
networks:
  default:
    name: fd_default
volumes:
  tei_data:
    name: fd_tei_data

```

## Proof Result

```json
{
  "apply": {
    "cmd": "docker compose up -d tei",
    "returncode": 0,
    "stdout": "",
    "stderr": " Container fd_tei  Recreate\n Container fd_tei  Recreated\n Container fd_tei  Starting\n Container fd_tei  Started\n",
    "duration_ms": 1045.55
  },
  "polls": [
    {
      "at": "2026-06-14T11:41:12Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:17Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:22Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:27Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:32Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:37Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:42Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:47Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:53Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:41:58Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:03Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:08Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:13Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:18Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:23Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:28Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:33Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:38Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:43Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:48Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:53Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:42:58Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:03Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:08Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:13Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:18Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:23Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:28Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:33Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:38Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:43Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:48Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:53Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:43:58Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:03Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:08Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:13Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:18Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:23Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:28Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:33Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:38Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:43Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:49Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:54Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:44:59Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:04Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:09Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:14Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:19Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:24Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:29Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:34Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:39Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:44Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:49Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:54Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:45:59Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:04Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:09Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:14Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:19Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:24Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:29Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:34Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:39Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:44Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:49Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:54Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:46:59Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:04Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:09Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:14Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:19Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:24Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:29Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:35Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:40Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:45Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:50Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:47:55Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:00Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:05Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:10Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:15Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:20Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:25Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:30Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:35Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:40Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:45Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:50Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:48:55Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:00Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:05Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:10Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:15Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:20Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:25Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:30Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:35Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:40Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:45Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:50Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:49:55Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:00Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:05Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:10Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:15Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:20Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:26Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:31Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:36Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:41Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:46Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:51Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:50:56Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:01Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:06Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:11Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:16Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:21Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:26Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:31Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:36Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:41Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:46Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:51Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:51:56Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:01Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:06Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:11Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:16Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:21Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:26Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:31Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:36Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:41Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:46Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:51Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:52:56Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:01Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:06Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:11Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:17Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:22Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:27Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:32Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:37Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:42Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:47Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:52Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:53:57Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:02Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:07Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:12Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:17Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:22Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:27Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:32Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:37Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:42Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:47Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:52Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:54:57Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:02Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:07Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:12Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:17Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:22Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:27Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:32Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:37Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:42Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:47Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:52Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:55:57Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:56:03Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    },
    {
      "at": "2026-06-14T11:56:08Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "unhealthy",
        "started_at": "2026-06-14T11:41:12.548077677Z",
        "safe_env": [
          "HF_HUB_OFFLINE=1",
          "HF_HOME=/data",
          "HUGGINGFACE_HUB_CACHE=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "deepvk/USER-bge-m3"
        ]
      }
    }
  ],
  "healthy_at": null,
  "outcome": "timeout",
  "logs_path": "/tmp/m045-tei-offline-startup.log",
  "post_state": [
    {
      "name": "fd_tei",
      "status": "running",
      "health": "unhealthy",
      "started_at": "2026-06-14T11:41:12.548077677Z",
      "safe_env": [
        "HF_HUB_OFFLINE=1",
        "HF_HOME=/data",
        "HUGGINGFACE_HUB_CACHE=/data",
        "PORT=80"
      ],
      "cmd": [
        "--model-id",
        "deepvk/USER-bge-m3"
      ]
    },
    {
      "name": "fd_api",
      "status": "running",
      "health": "healthy",
      "started_at": "2026-06-14T08:30:46.620802937Z",
      "safe_env": [
        "MODEL_ID=deepvk/USER-bge-m3",
        "PORT=8000",
        "TEI_URL=http://tei:80"
      ],
      "cmd": null
    },
    {
      "name": "fd_redis",
      "status": "running",
      "health": "healthy",
      "started_at": "2026-05-19T18:08:09.5269274Z",
      "safe_env": [],
      "cmd": [
        "redis-server",
        "--maxmemory",
        "2gb",
        "--maxmemory-policy",
        "allkeys-lru",
        "--save",
        "300",
        "1",
        "--appendonly",
        "no",
        "--protected-mode",
        "no"
      ]
    }
  ],
  "notable_log_lines": [
    "fd_tei  | \u001b[2m2026-06-14T11:41:12.752743Z\u001b[0m \u001b[32m INFO\u001b[0m \u001b[2mtext_embeddings_router\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/main.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m216:\u001b[0m Args { model_id: \"dee***/****-**e-m3\", revision: None, tokenization_workers: None, dtype: None, served_model_name: None, pooling: None, max_concurrent_requests: 512, max_batch_tokens: 16384, max_batch_requests: None, max_client_batch_size: 32, auto_truncate: true, default_prompt_name: None, default_prompt: None, dense_path: None, hf_api_token: None, hf_token: None, hostname: \"90ed1597262b\", port: 80, uds_path: \"/tmp/text-embeddings-inference-server\", huggingface_hub_cache: Some(\"/data\"), payload_limit: 2000000, api_key: None, json_output: false, disable_spans: false, otlp_endpoint: None, otlp_service_name: \"text-embeddings-inference.server\", prometheus_port: 9000, cors_allow_origin: None }",
    "fd_tei  | \u001b[2m2026-06-14T11:41:12.950845Z\u001b[0m \u001b[32m INFO\u001b[0m \u001b[2mtext_embeddings_router\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m271:\u001b[0m Starting model backend",
    "fd_tei  | \u001b[2m2026-06-14T11:50:14.871508Z\u001b[0m \u001b[32m INFO\u001b[0m \u001b[2mtext_embeddings_backend\u001b[0m\u001b[2m:\u001b[0m \u001b[2mbackends/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m706:\u001b[0m Downloading `onnx/model.onnx`"
  ]
}
```

## Notable TEI Log Lines

```text
fd_tei  | [2m2026-06-14T11:41:12.752743Z[0m [32m INFO[0m [2mtext_embeddings_router[0m[2m:[0m [2mrouter/src/main.rs[0m[2m:[0m[2m216:[0m Args { model_id: "dee***/****-**e-m3", revision: None, tokenization_workers: None, dtype: None, served_model_name: None, pooling: None, max_concurrent_requests: 512, max_batch_tokens: 16384, max_batch_requests: None, max_client_batch_size: 32, auto_truncate: true, default_prompt_name: None, default_prompt: None, dense_path: None, hf_api_token: None, hf_token: None, hostname: "90ed1597262b", port: 80, uds_path: "/tmp/text-embeddings-inference-server", huggingface_hub_cache: Some("/data"), payload_limit: 2000000, api_key: None, json_output: false, disable_spans: false, otlp_endpoint: None, otlp_service_name: "text-embeddings-inference.server", prometheus_port: 9000, cors_allow_origin: None }
fd_tei  | [2m2026-06-14T11:41:12.950845Z[0m [32m INFO[0m [2mtext_embeddings_router[0m[2m:[0m [2mrouter/src/lib.rs[0m[2m:[0m[2m271:[0m Starting model backend
fd_tei  | [2m2026-06-14T11:50:14.871508Z[0m [32m INFO[0m [2mtext_embeddings_backend[0m[2m:[0m [2mbackends/src/lib.rs[0m[2m:[0m[2m706:[0m Downloading `onnx/model.onnx`
```
