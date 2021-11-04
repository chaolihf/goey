{ pkgs ? import <nixpkgs> { }, platform ? "gtk", enableLint ? false }:

let
  platformInputs = with pkgs; {
    "gtk" = [ gtk3 pkg-config ];
    "js" = [ chromium wasmbrowsertest ];
    "windows" = [ wine64 fontconfig ];
  };
  platformEnvironment = {
    "gtk" = { };
    "js" = {
      GOOS = "js";
      GOARCH = "wasm";
    };
    "windows" = {
      GOOS = "windows";
      WINEDEBUG = "err+all,fixme-all";
    };
  };
  platformShellHook = {
    "gtk" = "";
    "js" = "chromium --version;"
      # See https://github.com/golang/go/issues/49011
      + "unset buildInputs HOST_PATH nobuildPhase builder shellHook phases; unset $(env | grep NIX_ | awk -F= '{print $1}');";
    "windows" = "wine64 --version;";
  };
  wasmbrowsertest = import ./wasmbrowsertest.nix { inherit pkgs; };
  wine64 = import ./wine64.nix { inherit pkgs; };
in pkgs.mkShell ({
  # No dependencies beyond stdlib.
  buildInputs = [ pkgs.go ] ++ platformInputs.${platform}
    ++ (if enableLint then [ pkgs.golangci-lint ] else [ ]);

  GOROOT = pkgs.go + "/share/go";

  # Make sure we don't pick up the users' GOPATH.
  # Advertise the current version of go when shell starts.
  shellHook = "unset GOPATH; go version;" + platformShellHook.${platform};
} // platformEnvironment.${platform})
