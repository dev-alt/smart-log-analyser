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
backups/, temp/, cache/
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
- ❌ SSL certificates or credential files
- ❌ Real log files with production data
- ❌ Environment files with real values (.env)
- ❌ Database connection strings
- ❌ Any file containing "password", "secret", "key", "token"

### Always Use:
- ✅ Example files with placeholder values
- ✅ Template configurations (servers.json.example)
- ✅ Environment variable references
- ✅ Localhost/example.com for examples
- ✅ Dummy/test credentials in documentation
- ✅ Clear security warnings in README

### Always Exclude in .gitignore:
- ✅ Real configuration files (servers.json, .env*)
- ✅ SSH key files (id_*, *.pem, *.key, *.crt)
- ✅ Download directories and log files
- ✅ Backup and temporary files
- ✅ IDE-specific files with potential secrets
- ✅ Any directory that might contain real data

---

## Documentation Standards

### README.md Requirements:
- **Feature documentation** for every new capability
- **Usage examples** with safe placeholder values
- **Security warnings** for sensitive features
- **Installation and setup** instructions
- **Command line options** documentation
- **Security notes** section

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
- [ ] No production log files or data
- [ ] .gitignore updated for new sensitive patterns
- [ ] Documentation updated for new features
- [ ] Security warnings added where appropriate
- [ ] Example files use placeholder values
- [ ] .development_log.md updated with session details

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