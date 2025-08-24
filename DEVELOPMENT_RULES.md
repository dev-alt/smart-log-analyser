# Smart Log Analyser Development Rules

**Established**: Session 6  
**Purpose**: Ensure consistent, secure, and well-documented development practices

---

## Mandatory Workflow for ALL Future Development

### 1. Documentation First 📚
- **Always update README.md** for any new features
- **Update relevant documentation** files (this file, examples, etc.)
- **Maintain .development_log.md** with session details including:
  - User instructions/requests
  - Implementation steps taken by Claude
  - Files created/modified and reasoning
  - Security considerations and decisions

### 2. Security Review 🔐
**Before every commit, verify:**
- ✅ Check all files for sensitive data (passwords, keys, IPs)
- ✅ Verify .gitignore excludes new sensitive patterns  
- ✅ Review for SSH keys, passwords, API keys, server details
- ✅ Use example/template files for sensitive configurations
- ✅ Never commit real credentials or production data

**Security Exclusions Checklist:**
```
# SSH Keys and Certificates
*.pem, *.key, *.crt, *.p12, *.pfx
id_*, *_rsa*, *_ed25519*, *_ecdsa*

# Configuration Files
.env*, config.json, servers.json (real configs)
*.conf (if contains credentials)

# Application Data
downloads/, logs/, *.log (real data)
output/, backups/, temp/, cache/
*.csv, *.json (analysis reports)
```

### 3. Development Session Tracking 📝
**For every development session, document:**
- User's exact instructions/requests
- Claude's interpretation and approach
- Step-by-step implementation details
- Files created/modified with explanations
- Testing performed
- Security considerations
- Any issues encountered and resolutions

### 4. Git Workflow 🚀
**Standard sequence for every development session:**
```bash
# 1. Stage all changes
git add .

# 2. Commit with descriptive message
git commit -m "Descriptive commit message with feature summary

- Key changes made
- Files affected
- Security considerations

🤖 Generated with [Claude Code](https://claude.ai/code)
Co-Authored-By: Claude <noreply@anthropic.com>"

# 3. Push to GitHub  
git push

# 4. Verify no sensitive data in commit history
git log --oneline -3
```

### 5. Testing & Validation ✅
**Before every commit:**
- ✅ Test new features work correctly
- ✅ Verify help commands display properly
- ✅ Ensure existing functionality still works
- ✅ Test that .gitignore exclusions work
- ✅ Verify no sensitive data is staged for commit

---

## Security Standards

### Never Commit:
- ❌ SSH private keys (id_rsa, id_ed25519, *.pem, *.key)
- ❌ Passwords or API keys in any format
- ❌ Real server IPs, hostnames, or connection details
- ❌ **REAL DOMAIN NAMES** - Never commit actual website domains
- ❌ **SPECIFIC IP ADDRESSES** - Never commit production server IPs
- ❌ SSL certificates or credential files
- ❌ Real log files with production data
- ❌ Environment files with real values (.env)
- ❌ Database connection strings
- ❌ Any file containing "password", "secret", "key", "token"
- ❌ **CLIENT/CUSTOMER IDENTIFIERS** - Never expose client-specific information

### Always Use:
- ✅ Example files with placeholder values
- ✅ Template configurations (servers.json.example)
- ✅ Environment variable references
- ✅ **GENERIC DOMAINS**: example.com, test.com, sample-site.com
- ✅ **PRIVATE IP RANGES**: 192.168.1.100, 10.0.0.50, 127.0.0.1
- ✅ **PLACEHOLDER PATHS**: server-logs/access.log, /path/to/logs/
- ✅ Dummy/test credentials in documentation
- ✅ Clear security warnings in README

### Always Exclude in .gitignore:
- ✅ Real configuration files (servers.json, .env*)
- ✅ SSH key files (id_*, *.pem, *.key, *.crt)
- ✅ Download directories and log files (downloads/, *.log)
- ✅ Output files and reports (output/, *.csv, *.json, detailed_report.*, summary.*)
- ✅ Backup and temporary files
- ✅ IDE-specific files with potential secrets
- ✅ Any directory that might contain real data or sensitive analysis results

---

## Project Structure Standards

### Folder Organization:
```
smart-log-analyser/
├── config/          # Configuration files (future use, excluded from git if sensitive)
├── downloads/       # Downloaded log files (ALWAYS excluded from git)
├── output/          # Generated reports and exports (ALWAYS excluded from git)  
├── testdata/        # Sample/test log files (safe for git - no real data)
├── pkg/             # Go packages (included in git)
├── cmd/             # CLI commands (included in git)
├── scripts/         # Utility scripts (included in git, check for sensitive data)
└── docs/            # Additional documentation (included in git)
```

### Folder Security Rules:
- **config/**: May contain sensitive data when implemented - verify before commits
- **downloads/**: NEVER commit - contains real log files with potentially sensitive data
- **output/**: NEVER commit - contains analysis results that may expose sensitive information
- **testdata/**: Safe to commit - contains only sanitized sample data
- **scripts/**: Review carefully - may contain temporary sensitive data or credentials

### New Folder Guidelines:
- Any new folder that might contain real data must be added to .gitignore
- Document the purpose and security considerations in folder README.md files
- Use placeholder/example files for any configuration templates

---

## Documentation Standards

### README.md Requirements:
- **Feature documentation** for every new capability
- **Usage examples** with safe placeholder values
- **Security warnings** for sensitive features
- **Installation and setup** instructions
- **Command line options** documentation
- **Security notes** section
- **Interactive mode documentation** for menu-driven workflows
- **HTML report examples** with browser integration instructions

### Code Documentation:
- **Clear function comments** for exported functions
- **Security warnings** in code near sensitive operations
- **Example usage** in function documentation
- **Error handling** explanations

### .development_log.md Format:
```markdown
### Session X: [Title]
**User Request**: "[Exact user instruction]"

**Tasks Completed**:
1. ✅ **Task description**
   - Implementation details
   - Files affected
   - Security considerations

**Files Added/Modified**:
- filename.go - Description of changes
- README.md - Updated sections

**Security Review**:
- Verified no credentials committed
- Updated .gitignore for new patterns
- Added security warnings to docs
```

### Documentation Security Guidelines

**Critical Rule: ALL examples in documentation MUST use generic placeholders**

#### ✅ Safe Documentation Examples:
```bash
# Safe command examples
./smart-log-analyser analyse server-logs/access.log --details
./smart-log-analyser download --server example.com
./smart-log-analyser analyse /var/log/nginx/access.log --export-html=report.html

# Safe configuration examples  
"host": "your-server.com"
"host": "192.168.1.100" 
"log_path": "/path/to/logs/access.log"
```

#### ❌ Dangerous Documentation Examples:
```bash
# NEVER include real domains or IPs in documentation
./smart-log-analyser analyse downloads/realsite.com_logs.log  # ❌ Real domain
./smart-log-analyser download --server 123.456.78.90          # ❌ Real IP
```

#### Development Log Security:
- **User Requests**: Sanitize user quotes to remove sensitive data before documenting
- **Command Examples**: Always use placeholder paths and generic domains
- **Error Messages**: Redact any real server information from logged errors
- **File Paths**: Use generic paths like `server-logs/` instead of real timestamps/hostnames

---

## Emergency Security Procedures

### If Sensitive Data is Accidentally Committed:

1. **Immediate Actions:**
   ```bash
   # If not yet pushed
   git reset --soft HEAD~1  # Undo last commit
   git reset HEAD filename  # Unstage sensitive file
   
   # If already pushed (DANGEROUS - rewrites history)
   git revert <commit-hash>  # Safer option
   # OR contact GitHub support for sensitive data removal
   ```

2. **Rotate Compromised Credentials:**
   - Change any passwords/keys that were exposed
   - Update server configurations
   - Notify relevant stakeholders

3. **Review and Improve:**
   - Update .gitignore patterns
   - Review development workflow
   - Add additional security checks

---

## Compliance Checklist

**Before every commit, confirm:**
- [ ] No real passwords, API keys, or tokens
- [ ] No SSH private keys or certificates  
- [ ] No real server IPs or hostnames
- [ ] No production log files or data in downloads/
- [ ] No analysis reports with sensitive data in output/
- [ ] .gitignore updated for new sensitive patterns
- [ ] Documentation updated for new features
- [ ] Security warnings added where appropriate
- [ ] Example files use placeholder values
- [ ] .development_log.md updated with session details
- [ ] New folders properly documented and secured

**Before every push, confirm:**
- [ ] All tests pass
- [ ] Help commands work correctly  
- [ ] No sensitive data in git history
- [ ] README.md reflects current features
- [ ] Security notes are up to date

---

## Continuous Improvement

These rules should be:
- **Reviewed** after every major feature addition
- **Updated** when new security concerns arise
- **Enhanced** based on lessons learned
- **Followed consistently** by all contributors

Remember: **Security and Documentation are not optional - they are requirements.**

---

*These rules ensure the Smart Log Analyser project maintains high security standards and comprehensive documentation while enabling rapid development.*