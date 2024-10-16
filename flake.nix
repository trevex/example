{
  description = "example";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    let
      overlays = [
        # none, yet
      ];
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system overlays;
        };
      in
      rec {
        devShell = pkgs.mkShell rec {
          name = "example";

          buildInputs = with pkgs; [
            go
            golangci-lint
            delve
            gnumake
            kubernetes-helm
          ];
        };
      }
    );
}
