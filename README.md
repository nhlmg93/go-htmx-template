# go-htmx-template

This is a template project for building Go applications using HTMX for frontend interactivity.

## Project Structure

- **`pkg/`**: Contains core application packages (router, web, logging, env).
- **`templates/`**: HTML templates used for rendering views.
- **`main.go`**: Entry point of the application.

## Running the Application

This project uses `.air` for hot-reloading during development.

### Development

To run the application with auto-reloading:

```bash
make dev
```


### Build

To build the final executable:

```bash
make build
```

## Dependencies

- Go
- Air (for development server)
