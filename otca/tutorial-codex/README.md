# OTCA Practice Lab (Codex)

OTCA(OpenTelemetry Certified Associate) 학습을 위해 구성한 Go 기반 실습 환경입니다. 단일 Go 서비스가 OpenTelemetry SDK로 계측되어 OTLP(OTel Line Protocol)를 통해 Collector로 전달하고, 수집된 데이터는 Jaeger와 Prometheus/Grafana로 확인할 수 있습니다.

## 구성 요소
- **Go Demo App** (`cmd/otel-app`): HTTP API(루트, `/inventory`, `/healthz`)를 제공하며 trace와 metrics를 노출합니다.
- **OpenTelemetry Collector** (`config/otel-collector.yaml`): OTLP(gRPC/HTTP) 수신, Jaeger + Prometheus로 내보내기, 로깅 exporter 활성화.
- **Jaeger All-in-One**: trace 탐색 UI (`http://localhost:16686`).
- **Prometheus**: Collector가 노출한 metrics를 스크랩 (`http://localhost:9090`).
- **Grafana**: 미리 등록된 대시보드/데이터소스 (`http://localhost:3000`, 익명 로그인 허용).

## 빠른 시작
1. **사전 준비**
   - Docker & Docker Compose가 설치되어 있어야 합니다.
   - 선택: 로컬에서 코드를 빌드/검증하려면 Go 1.22 이상.
2. **의존성 정리** (네트워크가 허용될 때 1회 실행)
   ```bash
   cd tutorial-codex
   GOCACHE=$(pwd)/.cache/go-build go mod tidy
   ```
3. **스택 실행**
   ```bash
   docker compose up --build
   ```
   - 앱: `http://localhost:8080`
   - 헬스체크: `http://localhost:8080/healthz`
   - 인벤토리 API: `http://localhost:8080/inventory`
4. **부하 생성(선택)**
   ```bash
   ./scripts/load.sh 40
   ```
   `TARGET_URL` 환경 변수를 바꿔 다른 엔드포인트를 대상으로 할 수 있습니다.

## 실습 아이디어
- **Trace 확인**: Jaeger에서 `otca-lab-service` 서비스 trace를 조회하고 span attribute/오류를 분석합니다.
- **Metrics 분석**: Grafana의 `OTCA Inventory Overview` 대시보드를 열어 latency 분포와 요청 성공/실패율을 확인합니다.
- **Collector 파이프라인 수정**: `config/otel-collector.yaml`에 span processor 추가(예: attributes, tail_sampling) 후 재기동.
- **수출 대상 교체**: Prometheus 대신 OTLP/HTTP → 다른 백엔드(Tempo, VictoriaMetrics 등)로 변경해보기.
- **코드 계측 확장**: handler에 추가 span/metric을 정의하고 `go test` 또는 `go run`으로 검증.

## 트러블슈팅 팁
- 빌드/실행 중 모듈 다운로드 문제가 발생하면 프록시(GOPROXY) 설정이나 사설 레지스트리 사용을 고려하세요.
- Collector가 기동하지 않을 경우 `docker compose logs otel-collector`로 구성 오류를 확인합니다.
- Prometheus 메트릭이 안 보일 경우 `otel-collector:8888/metrics` 엔드포인트 접속으로 확인 후, scrape 타이밍을 조정합니다.

## 참고 링크
- [OpenTelemetry Go Quickstart](https://opentelemetry.io/docs/languages/go/getting-started/)
- [OTCA 시험 안내](https://www.cncf.io/training/certification/otca/)
