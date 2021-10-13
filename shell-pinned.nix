let
  pkgs = import (builtins.fetchTarball {
    name = "nixpkgs-unstable";
    url =
      "https://github.com/nixos/nixpkgs/archive/2cf9db0e3d45b9d00f16f2836cb1297bcadc475e.tar.gz";
    sha256 = "0sij1a5hlbigwcgx10dkw6mdbjva40wzz4scn0wchv7yyi9ph48l";
  }) { };
in import ./shell.nix { pkgs = pkgs; }
