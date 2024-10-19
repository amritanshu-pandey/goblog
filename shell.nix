let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-24.05";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };

in pkgs.mkShellNoCC {
  packages = with pkgs; [
    go
    gopls
    go-outline
    gopkgs
    gocode-gomod
    godef
    golint
    templ
  ];

  shellHook = ''
    go version
  '';
}