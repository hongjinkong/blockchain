# ⛓️ Blockchain from Scratch — Go & Python

> 라이브러리 없이 **해시 체인과 작업증명(Proof of Work)** 을 직접 구현해 보며 블록체인의 동작 원리를 익힌 학습 프로젝트입니다.
> 같은 개념을 **Go**와 **Python** 두 언어로 각각 구현하고, HTTP API로 블록을 채굴·조회할 수 있게 만들었습니다.

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)
![Python](https://img.shields.io/badge/Python-3776AB?logo=python&logoColor=white)
![Flask](https://img.shields.io/badge/Flask-000000?logo=flask&logoColor=white)
![SHA-256](https://img.shields.io/badge/Hash-SHA--256-555555)

---

## 🧩 구현한 핵심 개념

| 개념 | 구현 내용 |
| --- | --- |
| **블록 & 해시 체인** | 각 블록이 이전 블록의 해시(`PrevHash`)를 담아 체인으로 연결 — 하나라도 바뀌면 이후 해시가 모두 어긋남 |
| **SHA-256 해싱** | 블록 필드를 이어 붙인 문자열을 SHA-256으로 해싱 |
| **작업증명 (PoW)** | 해시가 `0000`으로 시작할 때까지 `Nonce`를 1씩 올리며 탐색 (난이도 4) |
| **블록 검증** | 인덱스 연속성 · 이전 해시 일치 · 해시 재계산 일치를 확인한 뒤에만 체인에 추가 |
| **트랜잭션** (Python) | 거래를 대기열에 모았다가 채굴 시 블록에 담고, 채굴 보상 트랜잭션 지급 |

---

## 🔁 동작 흐름

```
 GET /mine
     ↓
 마지막 블록 조회 → 새 블록 생성 (Index+1, PrevHash = 이전 블록 해시)
     ↓
 Nonce = 0, 1, 2, ... → SHA-256 → "0000..."으로 시작할 때까지 반복   ← 작업증명
     ↓
 블록 검증 (Index 연속 · PrevHash 일치 · 해시 재계산)
     ↓
 체인에 추가 → 전체 체인을 JSON으로 응답
```

---

## 🐹 Go 버전 — `main.go`

`gorilla/mux` 라우터 기반 웹 서버. 블록에 심박수(BPM) 데이터를 기록합니다.

```go
type Block struct {
    Index     int
    Timestamp string
    BPM       int
    Hash      string
    PrevHash  string
    Nonce     int
}
```

| Method | Endpoint | 설명 |
| --- | --- | --- |
| `GET` | `/` | 전체 블록체인 조회 |
| `GET` | `/mine?bpm=72` | BPM 값을 담은 새 블록 채굴 |

```bash
go mod download
go run main.go                      # http://localhost:8080
curl "http://localhost:8080/mine?bpm=72"
```

> `main_basic.go`는 PoW와 웹 서버를 붙이기 전의 최소 블록 구조 프로토타입입니다.
> 같은 `main` 패키지라 함께 빌드되지 않으므로 `go run main.go`처럼 파일을 지정해 실행합니다.

---

## 🐍 Python 버전 — `my_blockchain/blockchain.py`

Flask 기반. 트랜잭션 대기열과 채굴 보상까지 포함한 확장 버전입니다.

| Method | Endpoint | 설명 |
| --- | --- | --- |
| `GET` | `/mine` | PoW로 새 블록 채굴 (채굴 보상 트랜잭션 포함) |
| `POST` | `/transactions/new` | 트랜잭션 추가 (`sender`, `recipient`, `amount`) |
| `GET` | `/chain` | 전체 체인과 길이 조회 |

```bash
pip install flask
python my_blockchain/blockchain.py  # http://localhost:5000

curl -X POST http://localhost:5000/transactions/new \
     -H "Content-Type: application/json" \
     -d '{"sender": "A", "recipient": "B", "amount": 5}'
curl http://localhost:5000/mine
curl http://localhost:5000/chain
```

---

## 🔀 두 구현 비교

| | Go | Python |
| --- | --- | --- |
| 웹 프레임워크 | `net/http` + `gorilla/mux` | Flask |
| 블록 데이터 | BPM(심박수) 값 | 트랜잭션 목록 |
| PoW 방식 | 블록 해시 자체가 `0000`으로 시작 | 이전 proof와 새 proof를 이은 값의 해시가 `0000`으로 시작 |
| 해싱 대상 | 필드 문자열 연결 | 블록 딕셔너리를 정렬된 JSON으로 직렬화 |
| 포트 | 8080 | 5000 |

---

## 📂 프로젝트 구조

```
blockchain/
├── main.go                     # Go — PoW + HTTP API 블록체인
├── main_basic.go               # Go — 최소 블록 구조 프로토타입
├── go.mod / go.sum             # Go 모듈 (gorilla/mux)
└── my_blockchain/
    └── blockchain.py           # Python — 트랜잭션 + PoW + Flask API
```

---

## 📌 한계

학습용 단일 노드 구현으로, 블록체인은 메모리에만 저장되어 서버를 재시작하면 초기화됩니다.
노드 간 체인 동기화(합의 알고리즘), 전자서명 기반 거래 검증, 영속 저장소는 구현 범위에 포함하지 않았습니다.
