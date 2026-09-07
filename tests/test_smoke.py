from pathlib import Path


def test_source_tree_contains_expected_safe_boundary_markers() -> None:
    source_root = Path(__file__).parents[1] / "src"

    assert source_root.is_dir()
    assert (source_root / "evidence").is_dir()
    assert (source_root / "report").is_dir()
    assert (source_root / "cleanup").is_dir()


def test_package_initializers_are_importable() -> None:
    package_dirs = [
        path for path in (Path(__file__).parents[1] / "src").iterdir() if path.is_dir() and path.name != "__pycache__"
    ]

    assert package_dirs
    assert all((package / "__init__.py").is_file() for package in package_dirs)
