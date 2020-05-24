{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [ pkgs.go pkgs.gtk3 pkgs.pkg-config ];

  shellHook = ''
    export GOPATH="$(pwd)/.go"
  '';
}
