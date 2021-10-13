{ pkgs ? import <nixpkgs> { } }:
let
  wasmbrowsertest = pkgs.buildGoModule rec {
    pname = "wasmbrowsertest";
    version = "0.3.5";

    src = pkgs.fetchFromGitHub {
      owner = "agnivade";
      repo = "wasmbrowsertest";
      rev = "v${version}";
      sha256 = "1n4sj35nzmckdlwipsdf184mmc70k17drjvw2m745jk4b8sjsw2f";
    };

    vendorSha256 = "1rg8qk520716gz11ydxh5ngkg24ipa6cg7jm7g6dwgl3b7znx7sk";

    meta = with pkgs.lib; {
      description = "Run Go wasm tests easily in your browser.";
      homepage = "https://github.com/agnivade/wasmbrowsertest";
      license = licenses.mit;
      platforms = platforms.linux ++ platforms.darwin;
    };
  };
in pkgs.stdenv.mkDerivation {
  pname = "go_js_wasm_exec";
  version = wasmbrowsertest.version;

  propagatedNativeBuildInputs = [ wasmbrowsertest ];

  src = ./wasmbrowsertest.nix;

  dontUnpack = true;
  dontConfigure = true;
  dontBuild = true;
  installPhase = ''
    mkdir -p $out/bin
    ln -s ${wasmbrowsertest}/bin/wasmbrowsertest $out/bin/go_js_wasm_exec 
  '';
}
