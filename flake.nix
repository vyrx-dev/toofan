{
  description = "A minimal, lightning-fast typing TUI";

  inputs = {
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    inputs:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } (
      { ... }: {
        systems = [
          "x86_64-linux"
        ];

        perSystem = { self', pkgs, ... }: {
          devShells.default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [ go ];
          };

          packages.default = pkgs.buildGoModule {
            pname = "toofan";
            version = "2.4.1";
            src = ./.;

            vendorHash = "sha256-YSjJ8NOL97hXZLnfGYIjoKmARv+gWOsv+5qkl9konnA=";

            nativeBuildInputs = with pkgs; [ go ];

            # `buildGoModule` was complaining about "inconsistent vendoring".
            # May need to rerun `go mod vendor` on the upstream repo.
            # This has nix stip the vendor dir and let's it vendor itself.
            # This fixes `nix run` and `nix build`.
            postPatch = ''
              rm -rf vendor
            '';

            meta = {
              description = "A minimal, lightning-fast typing TUI";
              mainProgram = "toofan";
            };
          };

          apps.default = {
            type = "app";
            program = "${self'.packages.default}/bin/toofan";
          };
        };
      }
    );
}
