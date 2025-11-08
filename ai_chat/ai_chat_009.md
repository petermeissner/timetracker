# AI Chat 009 - Command-Line Parameters Implementation

**Date**: November 8, 2025  
**Topic**: Adding configurable database file and port parameters  
**Status**: Completed ✅

## Overview

Enhanced the timesheet application with command-line parameter support to allow users to customize the database file path and server port, with additional environment variable support for flexible deployment scenarios.

## User Request

**Request**: "add parameters to the application so that i can modify the following:
- databse file used
- port used"

## Implementation

### Core Changes Made

#### 1. Enhanced main.go with Flag Support
- Added `flag` package import for command-line argument parsing
- Added `os` package import for environment variable support
- Implemented `getEnvOrDefault()` helper function

#### 2. Command-Line Parameters Added
- **`-port`**: Server port configuration (default: "8080")
- **`-db`**: SQLite database file path (default: "./timesheet.db")
- **`-help`**: Comprehensive usage information

#### 3. Environment Variable Support
- **`PORT`**: Alternative to -port flag
- **`DB_PATH`**: Alternative to -db flag
- **Precedence**: Command-line flags override environment variables

### Technical Implementation

#### Updated main.go Structure
```go
import (
    "database/sql"
    "embed"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    
    timesheet "timesheet/go"
    _ "modernc.org/sqlite"
)

// Helper function for environment variable fallbacks
func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func main() {
    // Define command-line flags with environment variable fallbacks
    var dbPath = flag.String("db", getEnvOrDefault("DB_PATH", "./timesheet.db"), "Path to the SQLite database file")
    var port = flag.String("port", getEnvOrDefault("PORT", "8080"), "Port to run the server on")
    var help = flag.Bool("help", false, "Show usage information")
    
    // Parse and handle flags
    flag.Parse()
    
    // Comprehensive help system
    if *help {
        // ... detailed help output with examples
    }
    
    // Use configured parameters
    serverAddr := "127.0.0.1:" + *port
    fmt.Printf("Server starting on http://localhost:%s\n", *port)
    fmt.Printf("Database file: %s\n", *dbPath)
    
    // ... rest of application logic
}
```

#### Key Features Implemented
1. **Flexible Configuration**: Multiple ways to configure the application
2. **Environment Integration**: Supports both development and production environments
3. **Comprehensive Help**: Detailed usage information with practical examples
4. **Backward Compatibility**: Maintains existing default behavior
5. **Parameter Validation**: Built-in Go flag validation and error handling

## Usage Examples

### Command-Line Flags
```powershell
# Default settings
.\timesheet.exe

# Custom port
.\timesheet.exe -port 8081

# Custom database file
.\timesheet.exe -db ./production.db

# Custom port and database
.\timesheet.exe -port 8081 -db ./production.db

# Show comprehensive help
.\timesheet.exe -help
```

### Environment Variables
```powershell
# PowerShell syntax
$env:PORT="8081"; .\timesheet.exe
$env:DB_PATH="./production.db"; .\timesheet.exe

# Combined environment variables
$env:PORT="8081"; $env:DB_PATH="./production.db"; .\timesheet.exe

# Command-line flags override environment variables
$env:PORT="8081"; .\timesheet.exe -port 8082  # Uses port 8082
```

### Development Workflows
```powershell
# Development with go run
go run main.go -port 8081 -db ./dev.db

# Production deployment
$env:PORT="80"; $env:DB_PATH="/var/lib/timesheet/production.db"; .\timesheet.exe

# Testing with isolated database
.\timesheet.exe -db ./test.db -port 8999
```

## Help System

The application now provides comprehensive built-in help:

```
Usage: timesheet.exe [options]

Options:
  -db string
        Path to the SQLite database file (default "./timesheet.db")
  -help
        Show usage information
  -port string
        Port to run the server on (default "8080")

Environment Variables:
  PORT      Port to run the server on (overridden by -port flag)
  DB_PATH   Path to the SQLite database file (overridden by -db flag)

Examples:
  timesheet.exe                              # Use default database and port
  timesheet.exe -port 8081                   # Use port 8081
  timesheet.exe -db ./custom.db              # Use custom database file
  timesheet.exe -db ./custom.db -port 8081   # Use custom database and port

  # Using environment variables:
  PORT=8081 timesheet.exe                    # Use port 8081
  DB_PATH=./custom.db timesheet.exe          # Use custom database file
```

## Testing and Validation

### Functionality Verified
✅ **Application builds successfully** without compilation errors  
✅ **Help command** displays comprehensive usage information  
✅ **Custom ports** work correctly (tested 8081, 8082, 8083, 8084, 8085)  
✅ **Custom database files** are created and used properly  
✅ **Environment variables** are respected for both PORT and DB_PATH  
✅ **Command-line flags override** environment variables correctly  
✅ **Database migrations** work with custom database paths  
✅ **Server startup** displays configured port and database file  

### Test Scenarios Executed
1. **Default Configuration**: `.\timesheet.exe`
2. **Custom Port**: `.\timesheet.exe -port 8081`
3. **Custom Database**: `.\timesheet.exe -db ./test.db`
4. **Combined Parameters**: `.\timesheet.exe -db ./test.db -port 8083`
5. **Environment Variables**: `$env:PORT="8084"; .\timesheet.exe`
6. **Flag Override**: `$env:PORT="8084"; .\timesheet.exe -port 8085`
7. **Help Display**: `.\timesheet.exe -help`

## Documentation Updates

### README.md Enhancements
- Added comprehensive **Command-Line Parameters** section
- Documented both **command-line flags** and **environment variables**
- Provided **practical examples** for different scenarios
- Explained **precedence rules** (flags override environment variables)
- Updated **installation and setup** instructions
- Enhanced **built executable usage** documentation

### Key Documentation Sections Added
1. **Command-Line Flags**: Complete parameter reference
2. **Environment Variables**: Alternative configuration method
3. **Examples**: Practical usage scenarios
4. **Precedence**: Clear explanation of override behavior
5. **Development Workflows**: Examples for different use cases
6. **Production Deployment**: Environment variable patterns

## Benefits Achieved

### For Developers
- **Flexible Development**: Easy port switching during development
- **Isolated Testing**: Separate database files for different test scenarios
- **Environment Consistency**: Same binary works across environments

### For Deployment
- **Production Ready**: Environment variable configuration for containers/services
- **Multi-Instance**: Multiple instances with different ports and databases
- **Configuration Management**: External configuration without code changes

### For Users
- **Simple Defaults**: Works out-of-the-box with sensible defaults
- **Easy Customization**: Simple parameter changes for different needs
- **Self-Documenting**: Built-in help provides all necessary information

## Architecture Benefits

### Code Quality
- **Clean Separation**: Configuration logic separated from business logic
- **Standard Patterns**: Uses Go's standard `flag` package
- **Maintainable**: Clear parameter definitions and help text
- **Extensible**: Easy to add more parameters in the future

### Deployment Flexibility
- **Container Ready**: Environment variable support for containerization
- **Multi-Environment**: Same binary for development, staging, production
- **Service Integration**: Easy integration with process managers and orchestrators

## Files Modified

### Core Application
- **`main.go`**: Complete parameter system implementation
  - Added flag parsing and environment variable support
  - Implemented comprehensive help system
  - Enhanced server startup with configuration display

### Documentation
- **`README.md`**: Complete documentation of new parameter system
  - Command-line flags and environment variables
  - Examples and usage patterns
  - Integration with existing documentation

## Future Extensibility

The parameter system is designed for easy extension:
- **Additional Parameters**: Simple to add new configuration options
- **Configuration Files**: Can be extended to support config file loading
- **Validation**: Framework in place for parameter validation
- **Environment Integration**: Ready for advanced deployment scenarios

## Impact

This enhancement significantly improves the application's:
- **Deployment Flexibility**: Multiple configuration methods
- **Development Experience**: Easy testing with different configurations
- **Production Readiness**: Professional parameter handling
- **User Experience**: Self-documenting with comprehensive help

The timesheet application now supports professional deployment patterns while maintaining ease of use for development and simple deployments.