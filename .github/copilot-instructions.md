# Promo Bot - Telegram Bot with Monobank Payment Integration

Promo Bot is a Python-based Telegram bot that handles ticket sales through Monobank payment integration and logs data to Baserow. It provides webhook endpoints for payment notifications and can be deployed locally, via Docker, or on Fly.io.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Prerequisites
- Python 3.11+ (Python 3.12 works)
- Docker (for containerized deployment)
- Internet connectivity for package installation
- Bot token from Telegram BotFather
- Monobank webhook configuration (for production)

### Bootstrap and Setup
1. **Create virtual environment:**
   ```bash
   python3 -m venv .venv
   source .venv/bin/activate
   ```

2. **Install dependencies:**
   ```bash
   pip install --upgrade pip
   pip install -r requirements.txt
   ```
   **TIMING**: Installation takes 30-60 seconds normally. NEVER CANCEL. Set timeout to 120+ seconds.
   
   **KNOWN ISSUE**: Original requirements.txt had dependency conflicts (aiogram==2.25.1 requires aiohttp<3.9.0 but aiohttp==3.9.5 was specified). This has been fixed to use compatible versions.

3. **Environment configuration:**
   ```bash
   cp .env.example .env
   # Edit .env with actual values:
   # BOT_TOKEN=your_telegram_bot_token
   # BASEROW_TOKEN=optional_baserow_token
   # BASEROW_TABLE_ID=optional_table_id
   # MONO_WEBHOOK_SECRET=optional_webhook_secret
   # TICKET_PRICE=10000
   ```

### Running the Application

#### Local Development
```bash
# After completing bootstrap steps above
source .venv/bin/activate
python main.py
```
**TIMING**: Application starts in 2-5 seconds. NEVER CANCEL. Set timeout to 30+ seconds.

The application runs both:
- Telegram bot (polling mode)
- HTTP webhook server on localhost:8080

#### Docker Deployment
```bash
# Build image
docker build -t promo-bot .
```
**TIMING**: Docker build takes 2-5 minutes depending on network speed. NEVER CANCEL. Set timeout to 600+ seconds.

```bash
# Run with environment file
docker run --rm -it -p 8080:8080 --env-file .env promo-bot
```

#### Fly.io Deployment
```bash
# Initial setup
fly launch --no-deploy
fly deploy
```
**TIMING**: Deployment takes 2-3 minutes. NEVER CANCEL. Set timeout to 300+ seconds.

```bash
# Set secrets
fly secrets set BOT_TOKEN=your_token MONO_WEBHOOK_SECRET=your_secret TICKET_PRICE=10000
fly secrets set BASEROW_TOKEN=your_token BASEROW_TABLE_ID=your_table_id
```

## Validation

### Health Check
```bash
curl http://localhost:8080/health
```
Should return "ok" when application is running.

### Application Testing Scenarios
**CRITICAL**: After making changes, ALWAYS test these scenarios:

1. **Basic Bot Functionality:**
   - Start the bot and send `/start` command
   - Verify "Купить билет" button appears
   - Click button and verify payment link generation

2. **Webhook Endpoint:**
   ```bash
   curl -X POST http://localhost:8080/mono-webhook \
     -H "Content-Type: application/json" \
     -d '{"amount":10000,"id":"test123","comment":"TEST1234"}'
   ```

3. **Configuration Validation:**
   - Verify bot starts without errors
   - Check logs for any missing configuration warnings
   - Confirm webhook server binds to correct port

### Manual Validation Requirements
- **ALWAYS** run through the complete user flow after making changes
- Test both successful and failed payment scenarios
- Verify error handling for invalid webhook data
- Check that pending tokens are properly managed

## Common Issues and Solutions

### Dependency Conflicts
- **Problem**: `aiogram==2.25.1` conflicts with `aiohttp==3.9.5`
- **Solution**: Use `aiohttp>=3.8.0,<3.9.0` (already fixed in requirements.txt)

### Network Connectivity
- **Problem**: Timeout errors during `pip install` or SSL certificate issues
- **Solutions**: 
  - Increase pip timeout: `pip install --timeout=300 -r requirements.txt`
  - Retry installation: pip often succeeds on second attempt
  - Use trusted hosts if SSL issues persist: `pip install --trusted-host pypi.org --trusted-host pypi.python.org --trusted-host files.pythonhosted.org -r requirements.txt`
  - For Docker builds, you may need to add SSL certificate handling or use base images with updated certificates

### Bot Token Issues
- **Problem**: Bot fails to start with authentication errors
- **Solution**: Verify `BOT_TOKEN` in `.env` file is correct and active

### Webhook Security
- **Problem**: Webhook receiving unauthorized requests
- **Solution**: Set `MONO_WEBHOOK_SECRET` and verify X-Mono-Secret header

## Code Structure

### Key Files
- `main.py` - Main application with bot handlers and webhook server
- `config.py` - Configuration management and environment variables
- `requirements.txt` - Python dependencies (fixed compatibility)
- `.env.example` - Environment variables template
- `Dockerfile` - Container configuration
- `fly.toml` - Fly.io deployment configuration

### Important Functions
- `generate_payment_token()` - Creates unique payment tracking tokens
- `add_ticket()` - Adds ticket data to Baserow and generates ticket number
- `mono_webhook()` - Handles Monobank payment notifications
- `start_cmd()` - Telegram bot /start command handler
- `buy_ticket()` - Handles ticket purchase button

### Database Integration
- Optional Baserow integration for ticket logging
- In-memory storage for pending payment tokens (production should use Redis)
- Ticket format: `T-{4-digit-random-number}`

## Development Workflow

### Making Changes
1. **Always** activate virtual environment: `source .venv/bin/activate`
2. Make code changes
3. Test locally: `python main.py`
4. Validate health endpoint: `curl http://localhost:8080/health`
5. Test bot functionality through Telegram
6. Build and test Docker image for production changes

### Testing
- No automated test suite exists - rely on manual testing
- Always test both Telegram bot and webhook endpoints
- Verify error handling and edge cases
- Test with both valid and invalid environment configurations

### Linting and Code Quality
- No linter configuration exists in repository
- Follow Python PEP 8 style guidelines
- Use meaningful variable names and add comments for complex logic

## Security Considerations

### Environment Variables
- **NEVER** commit actual secrets to repository
- Use `.env` file for local development
- Use Fly.io secrets for production deployment
- Bot token and webhook secrets are sensitive

### Webhook Security
- Implement HMAC signature verification for production
- Validate all incoming webhook data
- Use HTTPS in production (handled by Fly.io)

## Performance Notes

### Timing Expectations
- Virtual environment creation: 5-10 seconds
- Dependency installation: 30-120 seconds (network dependent)
- Application startup: 2-5 seconds
- Docker build: 2-5 minutes
- Fly.io deployment: 2-3 minutes

### Memory and Resources
- Application uses minimal resources (~50MB RAM)
- Single-threaded asyncio application
- Scales horizontally through multiple instances

## Troubleshooting

### Common Error Messages
- `BOT_TOKEN not set`: Missing or invalid bot token in environment
- `SSL certificate verification failed`: Network/firewall issues during pip install
- `Address already in use`: Port 8080 already occupied
- `Token not found`: Payment token expired or invalid

### Debug Mode
Add logging configuration for debugging:
```python
logging.basicConfig(level=logging.DEBUG)
```

### Port Conflicts
Change port in `.env` file:
```
PORT=8081
HOST=0.0.0.0
```

Remember: This is a production-ready bot that handles real payments. Always test thoroughly before deploying changes.

## Common Commands Reference

The following are outputs from frequently run commands. Reference them instead of viewing, searching, or running bash commands to save time.

### Repository Structure
```
ls -la /
.env.example          # Environment variables template  
.git/                 # Git repository data
.gitignore           # Git ignore patterns
Dockerfile           # Container configuration
README.md            # Project documentation (in Russian)
config.py            # Configuration and environment management
fly.toml             # Fly.io deployment configuration  
main.py              # Main application logic
requirements.txt     # Python dependencies (fixed compatibility)
```

### Python Files Compilation Test
```bash
python3 -m py_compile main.py config.py
# Should exit with code 0 - syntax is valid
```

### Sample Environment Configuration
```bash
cat .env.example
BOT_TOKEN=
BASEROW_TOKEN=
BASEROW_TABLE_ID=
MONO_WEBHOOK_SECRET=
RULES_VERSION=2025-09-05
TICKET_PRICE=10000
HOST=0.0.0.0
PORT=8080
```

### Requirements.txt Contents
```bash
cat requirements.txt
aiogram==2.25.1
aiohttp>=3.8.0,<3.9.0
python-dotenv==1.0.1
```

### Docker Health Check
The Dockerfile includes a health check:
```bash
python -c "import socket; s=socket.socket(); s.settimeout(2); s.connect(('127.0.0.1', 8080)); s.close()"
```