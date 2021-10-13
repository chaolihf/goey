{ pkgs ? import <nixpkgs> { } }:
let wine64 = pkgs.wineWowPackages.stable;
in pkgs.stdenv.mkDerivation {
  name = "go_windows_amd64_exec";

  propagatedNativeBuildInputs = [ wine64 ];

  src = ''
    #!{pkgs.bash}/bin/bash
    ${wine64}/bin/wine64 "$@"
  '';

  dontUnpack = true;
  dontConfigure = true;
  dontBuild = true;
  installPhase = ''
    mkdir -p $out/bin
    echo "$src" > $out/bin/go_windows_amd64_exec 
    chmod +x $out/bin/go_windows_amd64_exec
  '';
}
