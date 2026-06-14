# M045 TEI Local Path Startup Proof

Captured start: 2026-06-14T12:11:28Z

## Preflight

User directed the local model path approach after the `HF_HUB_OFFLINE=1` proof failed. The failed offline proof left TEI unhealthy/starting and showed that TEI still entered `Downloading onnx/model.onnx`; therefore offline env is rejected as the primary mitigation.

Selected local snapshot path:

`/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae`

Rollback if local-path proof fails:

1. Restore TEI command to `--model-id deepvk/USER-bge-m3` in both `docker-compose.yaml` and `docker-compose.override.yaml`.
2. Run `docker compose up -d tei`.
3. Wait for TEI health or document external TEI startup limitation.

Preflight container state:

```json
[
  {
    "name": "fd_tei",
    "status": "running",
    "health": "unhealthy",
    "started_at": "2026-06-14T12:01:17.468103026Z",
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

Compose candidate:

```text
name: fd
services:
  tei:
    command:
      - --model-id
      - /data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae
      - --max-batch-tokens
      - "8192"
    container_name: fd_tei
    environment:
      HF_HOME: /data
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

Local snapshot directory check:

```text
returncode=0
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/1_Pooling/config.json -> ../../../blobs/553a16bda12e2a6d2bb35de78c6ea264b7856e6a
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/config.json -> ../../blobs/883b8e471271d3bc817c517fe7d045e41f0fabf7
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/config_sentence_transformers.json -> ../../blobs/00601a94dee609765fc5bbfcd52c515e943c3683
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/model.safetensors -> ../../blobs/e6aa9c8e51a60ff383186a2f28f658305ba4ad23d2fa24296607885458ef2733
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/modules.json -> ../../blobs/952a9b81c0bfd99800fabf352f69c7ccd46c5e43
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/sentence_bert_config.json -> ../../blobs/0140ba1eac83a3c9b857d64baba91969d988624b
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/tokenizer.json -> ../../blobs/f61d51849a3308dde1b7c5f0cbfcc375e7b04ffe
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/tokenizer_config.json -> ../../blobs/c4d67b0826048e19f6c193896c6f717470613c7c


```

Success criteria:

- New TEI container command uses the local snapshot path.
- TEI reaches Docker health `healthy` within the proof timeout.
- TEI logs do not show remote Hub ONNX download attempts as the blocking path.
- fd `/health`, fd `/ready`, fd `/v1/embeddings`, and direct TEI `/embeddings` pass.

## Proof Result

```json
{
  "started_at": "2026-06-14T12:12:03Z",
  "apply": {
    "cmd": "docker compose up -d tei",
    "returncode": 0,
    "stdout": "",
    "stderr": " Container fd_tei  Recreate\n Container fd_tei  Recreated\n Container fd_tei  Starting\n Container fd_tei  Started\n",
    "duration_ms": 11070.17
  },
  "polls": [
    {
      "at": "2026-06-14T12:12:14Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:19Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:24Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:29Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:34Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:39Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:44Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:49Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:54Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:12:59Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:05Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:10Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:15Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:20Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:25Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:30Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:35Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:40Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:45Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:50Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:13:55Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:00Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:05Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:10Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:15Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:20Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:25Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:30Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:35Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:40Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:45Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:50Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:14:55Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:15:00Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:15:05Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:15:10Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "starting",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    },
    {
      "at": "2026-06-14T12:15:15Z",
      "tei": {
        "name": "fd_tei",
        "status": "running",
        "health": "healthy",
        "started_at": "2026-06-14T12:12:14.364809361Z",
        "safe_env": [
          "HUGGINGFACE_HUB_CACHE=/data",
          "HF_HOME=/data",
          "PORT=80"
        ],
        "cmd": [
          "--model-id",
          "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
          "--max-batch-tokens",
          "8192"
        ]
      }
    }
  ],
  "healthy_at": "2026-06-14T12:15:15Z",
  "outcome": "healthy",
  "logs_path": "/tmp/m045-tei-local-path-startup.log",
  "post_state": [
    {
      "name": "fd_tei",
      "status": "running",
      "health": "healthy",
      "started_at": "2026-06-14T12:12:14.364809361Z",
      "safe_env": [
        "HUGGINGFACE_HUB_CACHE=/data",
        "HF_HOME=/data",
        "PORT=80"
      ],
      "cmd": [
        "--model-id",
        "/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae",
        "--max-batch-tokens",
        "8192"
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
    "fd_tei  | \u001b[2m2026-06-14T12:12:14.621566Z\u001b[0m \u001b[32m INFO\u001b[0m \u001b[2mtext_embeddings_router\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/main.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m216:\u001b[0m Args { model_id: \"/dat*/******--******--****-***-**/*********/*************************************9ae\", revision: None, tokenization_workers: None, dtype: None, served_model_name: None, pooling: None, max_concurrent_requests: 512, max_batch_tokens: 8192, max_batch_requests: None, max_client_batch_size: 32, auto_truncate: true, default_prompt_name: None, default_prompt: None, dense_path: None, hf_api_token: None, hf_token: None, hostname: \"83e965b0efbf\", port: 80, uds_path: \"/tmp/text-embeddings-inference-server\", huggingface_hub_cache: Some(\"/data\"), payload_limit: 2000000, api_key: None, json_output: false, disable_spans: false, otlp_endpoint: None, otlp_service_name: \"text-embeddings-inference.server\", prometheus_port: 9000, cors_allow_origin: None }",
    "fd_tei  | \u001b[2m2026-06-14T12:12:14.732833Z\u001b[0m \u001b[32m INFO\u001b[0m \u001b[2mtext_embeddings_router\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m271:\u001b[0m Starting model backend",
    "fd_tei  | \u001b[2m2026-06-14T12:12:14.733517Z\u001b[0m \u001b[31mERROR\u001b[0m \u001b[2mtext_embeddings_backend\u001b[0m\u001b[2m:\u001b[0m \u001b[2mbackends/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m412:\u001b[0m Could not start ORT backend: Could not start backend: File at `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/onnx/model.onnx` does not exist",
    "fd_tei  | \u001b[2m2026-06-14T12:12:14.734145Z\u001b[0m \u001b[33m WARN\u001b[0m \u001b[2mtext_embeddings_backend_candle\u001b[0m\u001b[2m:\u001b[0m \u001b[2mbackends/candle/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m229:\u001b[0m The `config.json` contains `hidden_act=gelu` and GeLU + tanh approximation will be used instead of exact GeLU (aka. GeLU erf), which might lead to subtle differences with Transformers or Sentence Transformers outputs which use exact GeLU when `hidden_act=gelu`, unless specified otherwise. GeLU + tanh is more efficient and more consistent across devices (e.g., cuBLASLt comes with fused GeLU + tanh), and will have negligible impact on the inference quality.",
    "fd_tei  | \u001b[2m2026-06-14T12:12:14.735297Z\u001b[0m \u001b[32m INFO\u001b[0m \u001b[2mtext_embeddings_backend_candle\u001b[0m\u001b[2m:\u001b[0m \u001b[2mbackends/candle/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m311:\u001b[0m Starting Bert model on Cpu",
    "fd_tei  | \u001b[2m2026-06-14T12:12:18.460470Z\u001b[0m \u001b[32m INFO\u001b[0m \u001b[2mtext_embeddings_router\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m289:\u001b[0m Warming up model",
    "fd_tei  | \u001b[2m2026-06-14T12:15:08.132221Z\u001b[0m \u001b[33m WARN\u001b[0m \u001b[2mtext_embeddings_router\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m303:\u001b[0m Backend does not support a batch size > 4",
    "fd_tei  | \u001b[2m2026-06-14T12:15:08.132252Z\u001b[0m \u001b[33m WARN\u001b[0m \u001b[2mtext_embeddings_router\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m304:\u001b[0m forcing `max_batch_requests=4`",
    "fd_tei  | \u001b[2m2026-06-14T12:15:08.132650Z\u001b[0m \u001b[33m WARN\u001b[0m \u001b[2mtext_embeddings_router\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/lib.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m354:\u001b[0m Invalid hostname, defaulting to 0.0.0.0",
    "fd_tei  | \u001b[2m2026-06-14T12:15:08.134310Z\u001b[0m \u001b[32m INFO\u001b[0m \u001b[2mtext_embeddings_router::http::server\u001b[0m\u001b[2m:\u001b[0m \u001b[2mrouter/src/http/server.rs\u001b[0m\u001b[2m:\u001b[0m\u001b[2m1881:\u001b[0m Ready"
  ],
  "smoke": {
    "fd_health": {
      "http_status": 200,
      "status": "ok",
      "runtime": {
        "backend": "tei",
        "model": "deepvk/USER-bge-m3",
        "dimensions": 1024,
        "production_default": true,
        "cache_namespace": "v2"
      }
    },
    "fd_ready": {
      "http_status": 200,
      "body": {
        "status": "ready",
        "time": "2026-06-14T12:15:16Z"
      }
    },
    "fd_embedding": {
      "http_status": 200,
      "latency_ms": 308.85,
      "model": "deepvk/USER-bge-m3",
      "embedding_len": 1024
    },
    "tei_embedding": {
      "http_status": 200,
      "latency_ms": 305.98,
      "embedding_len": 1024
    }
  }
}
```

## Notable TEI Log Lines

```text
fd_tei  | [2m2026-06-14T12:12:14.621566Z[0m [32m INFO[0m [2mtext_embeddings_router[0m[2m:[0m [2mrouter/src/main.rs[0m[2m:[0m[2m216:[0m Args { model_id: "/dat*/******--******--****-***-**/*********/*************************************9ae", revision: None, tokenization_workers: None, dtype: None, served_model_name: None, pooling: None, max_concurrent_requests: 512, max_batch_tokens: 8192, max_batch_requests: None, max_client_batch_size: 32, auto_truncate: true, default_prompt_name: None, default_prompt: None, dense_path: None, hf_api_token: None, hf_token: None, hostname: "83e965b0efbf", port: 80, uds_path: "/tmp/text-embeddings-inference-server", huggingface_hub_cache: Some("/data"), payload_limit: 2000000, api_key: None, json_output: false, disable_spans: false, otlp_endpoint: None, otlp_service_name: "text-embeddings-inference.server", prometheus_port: 9000, cors_allow_origin: None }
fd_tei  | [2m2026-06-14T12:12:14.732833Z[0m [32m INFO[0m [2mtext_embeddings_router[0m[2m:[0m [2mrouter/src/lib.rs[0m[2m:[0m[2m271:[0m Starting model backend
fd_tei  | [2m2026-06-14T12:12:14.733517Z[0m [31mERROR[0m [2mtext_embeddings_backend[0m[2m:[0m [2mbackends/src/lib.rs[0m[2m:[0m[2m412:[0m Could not start ORT backend: Could not start backend: File at `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/onnx/model.onnx` does not exist
fd_tei  | [2m2026-06-14T12:12:14.734145Z[0m [33m WARN[0m [2mtext_embeddings_backend_candle[0m[2m:[0m [2mbackends/candle/src/lib.rs[0m[2m:[0m[2m229:[0m The `config.json` contains `hidden_act=gelu` and GeLU + tanh approximation will be used instead of exact GeLU (aka. GeLU erf), which might lead to subtle differences with Transformers or Sentence Transformers outputs which use exact GeLU when `hidden_act=gelu`, unless specified otherwise. GeLU + tanh is more efficient and more consistent across devices (e.g., cuBLASLt comes with fused GeLU + tanh), and will have negligible impact on the inference quality.
fd_tei  | [2m2026-06-14T12:12:14.735297Z[0m [32m INFO[0m [2mtext_embeddings_backend_candle[0m[2m:[0m [2mbackends/candle/src/lib.rs[0m[2m:[0m[2m311:[0m Starting Bert model on Cpu
fd_tei  | [2m2026-06-14T12:12:18.460470Z[0m [32m INFO[0m [2mtext_embeddings_router[0m[2m:[0m [2mrouter/src/lib.rs[0m[2m:[0m[2m289:[0m Warming up model
fd_tei  | [2m2026-06-14T12:15:08.132221Z[0m [33m WARN[0m [2mtext_embeddings_router[0m[2m:[0m [2mrouter/src/lib.rs[0m[2m:[0m[2m303:[0m Backend does not support a batch size > 4
fd_tei  | [2m2026-06-14T12:15:08.132252Z[0m [33m WARN[0m [2mtext_embeddings_router[0m[2m:[0m [2mrouter/src/lib.rs[0m[2m:[0m[2m304:[0m forcing `max_batch_requests=4`
fd_tei  | [2m2026-06-14T12:15:08.132650Z[0m [33m WARN[0m [2mtext_embeddings_router[0m[2m:[0m [2mrouter/src/lib.rs[0m[2m:[0m[2m354:[0m Invalid hostname, defaulting to 0.0.0.0
fd_tei  | [2m2026-06-14T12:15:08.134310Z[0m [32m INFO[0m [2mtext_embeddings_router::http::server[0m[2m:[0m [2mrouter/src/http/server.rs[0m[2m:[0m[2m1881:[0m Ready
```
