def find_docker_compose_files(directory):
    """Recursively finds all `docker-compose.yml` files in a directory.

    Args:
    directory: The directory to search as a string.

    Returns:
    A list of paths to the found `docker-compose.yml` files.
    """
    files = []

    # List all files and directories in the given directory
    entries = str(local("find {} -type f -name docker-compose.yml".format(directory))).split("\n")

    for entry in entries:
        if entry:  # Ensure no empty strings are added
            files.append(entry)

    return files


# Print the found docker-compose.yml files (or use them further as needed)
docker_compose_files = find_docker_compose_files(".")
for file in docker_compose_files:
    print("Found docker-compose file:", file)
if docker_compose_files:
    docker_compose(docker_compose_files)
