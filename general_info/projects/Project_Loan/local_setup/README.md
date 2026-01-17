# Local Development Setup for Loan Validator System

This folder contains everything you need to run the complete Loan Validator System locally.

## 🚀 Quick Start

### Start Everything (One Command)

```bash
cd /home/dapo/daas/daascohort3/general_info/projects/Project_Loan/local_setup
./start-all.sh
```

This will:
1. ✅ Start PostgreSQL in Docker
2. ✅ Create both databases
3. ✅ Run database migrations
4. ✅ Start OpenTelemetry Collector
5. ✅ Start Government API
6. ✅ Start Loan Validator Portal

### Stop Everything (One Command)

```bash
cd /home/dapo/daas/daascohort3/general_info/projects/Project_Loan/local_setup
./teardown.sh
```

This will:
1. Stop all Go applications
2. Stop Docker containers (OTel Collector, PostgreSQL)
3. Clean up processes

---

## 📁 Files in This Directory

| File | Purpose |
|------|---------|
| `start-all.sh` | **Main startup script** - Runs everything |
| `teardown.sh` | **Cleanup script** - Stops everything |
| `docker-compose.local.yml` | Docker services (PostgreSQL + OTel Collector) |
| `setup-databases.sh` | Database initialization script |
| `start-government-api.sh` | Starts Government API in background |
| `start-loan-validator.sh` | Starts Loan Validator Portal in background |
| `README.md` | This file |

---

## 🔧 Manual Control (Advanced)

### Start Components Individually

```bash
# 1. Start Docker services (PostgreSQL + OTel Collector)
docker-compose -f docker-compose.local.yml up -d

# 2. Setup databases
./setup-databases.sh

# 3. Start Government API
./start-government-api.sh

# 4. Start Loan Validator Portal
./start-loan-validator.sh
```

### View Application Logs

```bash
# Government API logs
tail -f logs/government-api.log

# Loan Validator Portal logs
tail -f logs/loan-validator-portal.log

# OTel Collector logs
docker logs -f otel-collector

# PostgreSQL logs
docker logs -f loan_postgres
```

### Check Status

```bash
# Check what's running
docker ps

# Check Go processes
ps aux | grep "go run"

# Test endpoints
curl http://localhost:8082/health  # Government API
curl http://localhost:8080          # Loan Validator Portal
```

---

## 🌐 Access Points

| Service | URL | Description |
|---------|-----|-------------|
| **Loan Validator Portal** | http://localhost:8080 | Main application |
| **Government API** | http://localhost:8082 | Customer data API |
| **Datadog APM** | https://us5.datadoghq.com/apm/traces | View traces |
| **PostgreSQL** | localhost:5432 | Database (user: postgres, pass: postgres) |
| **OTel Collector** | localhost:4317 | OTLP gRPC endpoint |

---

## 🧪 Testing the Setup

### 1. Access the Application

Open your browser to: **http://localhost:8080**

### 2. Register a User

Use any name from the sample data:
- **Oladapo Babalola**
- **John Doe**
- **Jane Smith**
- **Mike Wilson**
- **Sarah Johnson**

### 3. Login and Validate Loan

After login, click "Click Here to Validate Your Loan"

### 4. View Traces in Datadog

Visit: https://us5.datadoghq.com/apm/traces

Filter by:
- Service: `loan-validator-portal`
- Environment: `development`

---

## 🐛 Troubleshooting

### Problem: Port Already in Use

```bash
# Check what's using port 8080
lsof -i :8080

# Check what's using port 8082
lsof -i :8082

# Kill process if needed
kill -9 <PID>
```

### Problem: Database Connection Failed

```bash
# Check if PostgreSQL is running
docker ps | grep loan_postgres

# Restart PostgreSQL
docker-compose -f docker-compose.local.yml restart postgres

# Check database exists
docker exec loan_postgres psql -U postgres -l
```

### Problem: No Traces in Datadog

```bash
# Check OTel Collector
docker logs otel-collector | grep -i error

# Verify API key
docker logs otel-collector | grep "API key"

# Check if collector is receiving data
docker logs otel-collector | tail -20
```

### Problem: Application Won't Start

```bash
# Check Go dependencies
cd ../government_api/backend && go mod download
cd ../loan_validator_portal/backend && go mod download

# Check for existing processes
ps aux | grep "go run"

# Run teardown and try again
./teardown.sh
sleep 2
./start-all.sh
```

---

## 📊 What's Running After Start

```
┌─────────────────────────────────────┐
│   Loan Validator Portal             │
│   Port: 8080 (Go Process)           │
│   PID: stored in logs/validator.pid │
└────────────┬────────────────────────┘
             │
             │ HTTP API Call
             │
             ▼
┌─────────────────────────────────────┐
│   Government API                     │
│   Port: 8082 (Go Process)           │
│   PID: stored in logs/gov-api.pid   │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│   OpenTelemetry Collector            │
│   Port: 4317 (Docker)               │
│   Forwards traces to Datadog US5    │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│   PostgreSQL 15                      │
│   Port: 5432 (Docker)               │
│   Databases:                        │
│   - government_loan_db              │
│   - loan_validator_db               │
└─────────────────────────────────────┘
```

---

## 🔐 Environment Variables

All environment variables are automatically set by the startup scripts:

### Government API
- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=government_loan_db`
- `PORT=8082`

### Loan Validator Portal
- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=loan_validator_db`
- `PORT=8080`
- `GOV_BANK_URL=http://localhost:8082`
- `OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317`
- `ENV=development`

---

## 📝 Sample Data

The database is pre-populated with sample customers:

| First Name | Last Name | Loan Amount | Status |
|------------|-----------|-------------|--------|
| Oladapo | Babalola | $100,000 | Approved |
| John | Doe | $50,000 | Approved |
| Jane | Smith | $75,000 | Pending |
| Mike | Wilson | $30,000 | Rejected |
| Sarah | Johnson | $60,000 | Approved |

Use these names to test the loan validation feature.

---

## 🎯 Success Checklist

After running `./start-all.sh`:

- [ ] ✅ PostgreSQL container running
- [ ] ✅ OTel Collector container running
- [ ] ✅ Government API responding on port 8082
- [ ] ✅ Loan Validator Portal responding on port 8080
- [ ] ✅ Can access http://localhost:8080 in browser
- [ ] ✅ Can register and login
- [ ] ✅ Can validate loan
- [ ] ✅ Traces appearing in Datadog

---

## 🆘 Need Help?

1. **Check logs**: All logs are in the `logs/` directory
2. **Run teardown**: `./teardown.sh` to clean up
3. **Restart fresh**: `./teardown.sh && ./start-all.sh`
4. **Check Docker**: `docker ps` to see running containers
5. **Check processes**: `ps aux | grep "go run"` to see Go apps

---

## 🚀 You're All Set!

Just run `./start-all.sh` and open http://localhost:8080 in your browser!

The system will automatically:
- Set up all databases
- Start all services
- Configure observability
- Be ready to use in ~10 seconds
