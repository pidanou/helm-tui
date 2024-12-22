#!/usr/bin/env sh

# Shamelessly copied from https://github.com/databus23/helm-diff

PROJECT_NAME="helm-tui"
PROJECT_GH="pidanou/$PROJECT_NAME"
export GREP_COLOR="never"

# Convert HELM_BIN and HELM_PLUGIN_DIR to unix if cygpath is
# available. This is the case when using MSYS2 or Cygwin
# on Windows where helm returns a Windows path but we
# need a Unix path

if command -v cygpath >/dev/null 2>&1; then
  HELM_BIN="$(cygpath -u "${HELM_BIN}")"
  HELM_PLUGIN_DIR="$(cygpath -u "${HELM_PLUGIN_DIR}")"
fi

[ -z "$HELM_BIN" ] && HELM_BIN=$(command -v helm)

[ -z "$HELM_HOME" ] && HELM_HOME=$(helm env | grep 'HELM_DATA_HOME' | cut -d '=' -f2 | tr -d '"')

mkdir -p "$HELM_HOME"

: "${HELM_PLUGIN_DIR:="$HELM_HOME/plugins/helm-tui"}"

if [ "$SKIP_BIN_INSTALL" = "1" ]; then
  echo "Skipping binary install"
  exit
fi

# which mode is the common installer script running in
SCRIPT_MODE="install"
if [ "$1" = "-u" ]; then
  SCRIPT_MODE="update"
fi

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
  aarch64) ARCH="arm64" ;;
  x86_64) ARCH="x86_64" ;;
  i386) ARCH="i386" ;;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(uname -s)

  case "$OS" in
  Windows_NT) OS='Windows' ;;
  # Msys support
  MSYS*) OS='Windows' ;;
  # Minimalist GNU for Windows
  MINGW*) OS='Windows' ;;
  CYGWIN*) OS='Windows' ;;
  Darwin) OS='Darwin' ;;
  Linux) OS='Linux' ;;
  esac
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  supported="Linux_x86_64\nLinux_arm64\nLinux_i386\nDarwin_x86_64\nDarwin_arm64\nWindows_arm64\nWindows_i386\nWindows_x86_64"
  if ! echo "${supported}" | grep -q "${OS}_${ARCH}"; then
    echo "No prebuild binary for ${OS}_${ARCH}."
    exit 1
  fi

  if
    ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1
  then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# getDownloadURL checks the latest available version.
getDownloadURL() {
  version=$(git -C "$HELM_PLUGIN_DIR" describe --tags --exact-match 2>/dev/null || :)
  if [ "$SCRIPT_MODE" = "install" ] && [ -n "$version" ]; then
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/download/$version/helm-tui_${OS}_${ARCH}.tar.gz"
  else
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/latest/download/helm-tui_${OS}_${ARCH}.tar.gz"
  fi
}

# Temporary dir
mkTempDir() {
  HELM_TMP="$(mktemp -d -t "${PROJECT_NAME}-XXXXXX")"
}
rmTempDir() {
  if [ -d "${HELM_TMP:-/tmp/helm-tui-tmp}" ]; then
    rm -rf "${HELM_TMP:-/tmp/helm-tui-tmp}"
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  PLUGIN_TMP_FILE="${HELM_TMP}/${PROJECT_NAME}.tar.gz"

  echo "Downloading $DOWNLOAD_URL"
  if
    command -v curl >/dev/null 2>&1
  then
    curl -sSf -L "$DOWNLOAD_URL" >"$PLUGIN_TMP_FILE"
  elif
    command -v wget >/dev/null 2>&1
  then
    wget -q -O - "$DOWNLOAD_URL" >"$PLUGIN_TMP_FILE"
  fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  tar xzf "$PLUGIN_TMP_FILE" -C "$HELM_TMP"
  HELM_TMP_BIN="$HELM_TMP/helm-tui"
  if [ "${OS}" = "windows" ]; then
    HELM_TMP_BIN="$HELM_TMP_BIN.exe"
  fi
  echo "Preparing to install into ${HELM_PLUGIN_DIR}"
  mkdir -p "${HELM_PLUGIN_DIR}/bin"
  cp "$HELM_TMP_BIN" "$HELM_PLUGIN_DIR/bin/helm-tui"
}

# exit_trap is executed if on exit (error or not).
exit_trap() {
  result=$?
  rmTempDir
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    printf '\tFor support, go to https://github.com/pidanou/helm-tui.\n'
  fi
  exit $result
}

# Execution

#Stop execution on any error
trap "exit_trap" EXIT
set -e
initArch
initOS
verifySupported
getDownloadURL
mkTempDir
downloadFile
installFile
