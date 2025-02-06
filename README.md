# Helm-tui

<img alt="Demo of Soramail" src="demos/overview.gif" width="1200" />

Helm-tui is a terminal-based UI application to manage your Helm releases, charts, repositories, and plugins with ease.

## Features

- Manage Helm releases effortlessly.
- Add, update, and remove Helm repositories.

## Requirements

- [Helm 3](https://helm.sh/docs/intro/install/)

### Optional

- [Go 1.22+](https://go.dev/doc/install)

## How to Use

1. Clone the repository:

   ```bash
   git clone https://github.com/pidanou/helm-tui.git
   cd helm-tui
   ```

2. Run the app directly using:
   ```bash
   go run .
   ```

## How to Install

### Install Helm-tui using `helm plugin install`:

```bash
helm plugin install https://github.com/pidanou/helm-tui
```

Once installed, you can run `helm tui` directly from your terminal.


### Install using Go:

```bash
go install https://github.com/pidanou/helm-tui@latest
```

Once installed, you can run `helm-tui` directly from your terminal.

## Contributing

Contributions are welcome! If you find bugs or have feature requests, feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
