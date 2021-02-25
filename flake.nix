{
  description = "A flake for building launchd_shim.";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    with import nixpkgs;
    let
      name = "launchd_shim";
      systems = [ "x86_64-darwin" ];
      forAllSystems = flake-utils.lib.eachSystem systems;
      src = self;
    in
    {
      overlay = self: super: {
        ${name} = super.buildGoModule {
          inherit name src;
          version = "2021-02-24";
          vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
        };
      };
    } // (
      forAllSystems (system:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [ self.overlay ];
          };

          package = pkgs.${name};
        in {
          packages.${name} = package;
          defaultPackage = package;
        }
      )
    );
}
