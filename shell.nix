{ pkgs ? import <nixpkgs> { }, platform ? "gtk", enableLint ? false }:

let
  platformInputs = with pkgs; {
    "gtk" = [ gtk3 pkg-config ];
    "js" = [ chromium wasmbrowsertest ];
    "windows" = [ wine64 ];
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
  wasmbrowsertest = import ./wasmbrowsertest.nix { inherit pkgs; };
  wine64 = import ./wine64.nix { inherit pkgs; };
in pkgs.mkShell ({
  # No dependencies beyond stdlib.
  buildInputs = [ pkgs.go ] ++ platformInputs.${platform}
    ++ (if enableLint then [ pkgs.golangci-lint ] else [ ]);

  GOROOT = pkgs.go + "/share/go";

  # Make sure we don't pick up the users' GOPATH.
  # Advertise the current version of go when shell starts.
  shellHook = "unset GOPATH; go version;"
    + (if platform == "windows" then "wine64 --version;" else "");
} // platformEnvironment.${platform})
