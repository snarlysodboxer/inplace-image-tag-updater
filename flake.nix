{
  description = "Update image tags in Kubernetes Specs without touching formatting";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "inplace-image-tag-updater";
            version = "0.0.3";

            src = ./.;

            vendorHash = null;

            meta = with pkgs.lib; {
              description = "Update image tags in Kubernetes Specs without touching formatting";
              homepage = "https://github.com/snarlysodboxer/inplace-image-tag-updater";
              license = licenses.asl20;
              maintainers = with maintainers; [ snarlysodboxer ];
              mainProgram = "inplace-image-tag-updater";
            };
          };
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/inplace-image-tag-updater";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
          ];
        };
      }
    );
}
