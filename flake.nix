{
  inputs.nixpkgs.url = "nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    let
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;
      supportedSystems = [ "x86_64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
      libFor = forAllSystems (system: import (nixpkgs + "/lib"));
      nixosLibFor = forAllSystems (system: import (nixpkgs + "/nixos/lib"));
    in flake-utils.lib.eachSystem supportedSystems (system: let 
      pkgs = import nixpkgs {
        inherit system;
      };
      lib = import (nixpkgs + "/lib") {
        inherit system;
      };
      nixosLib = import (nixpkgs + "/nixos/lib") {
        inherit system;
      };
      ldflags = pkgs: [];
      deps = pkgs: with pkgs; [ libwebp ];
    in rec {
      devShells = let pkgs = nixpkgsFor.${system}; in { default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_19
            git
          ] ++ (deps pkgs);
      }; };
      packages = let
        pkgs = nixpkgsFor.${system};
        lib = libFor.${system};
        common = {
          inherit version;
          src = ./.;
          ldflags = ldflags pkgs;
          tags = [ "nix" "sdnotify" ];
          #vendorSha256 = pkgs.lib.fakeSha256;
          # run ./base64-hex on the sha256- hash returned when using fakeSha256 to update
          vendorSha256 = "ce0bb86722d916d3da3ae5b9384af2615af550d7df2c7489425ed3c513da7784";
          buildInputs = deps pkgs;
        };
      in
      {
        default = pkgs.buildGoModule (common // {
          pname = "komorebi-proxy";
          postInstall = ''
            mv $out/bin/proxy $out/bin/komorebi-proxy
          '';
        });
      };
  });
}
