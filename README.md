# Alt Galaxy

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Alternate implementation of the
[ansible-galaxy](http://docs.ansible.com/ansible/galaxy.html) tool for
installing Ansible roles.

## Warning

Consider this software alpha quality. This project was created as a temporary
workaround rather than a permanent solution. The long term intention is to
contribute improvements to the `ansible-galaxy` tool to resolve the issues
there.

## Objective

On projects with many roles `ansible-galaxy` sufferers from intermittent HTTP
failures when installing them. This project implements a limited subset of the
`ansible-galaxy` features to allow you to install roles more reliably.

If you've seen the following error you may find this project useful:

```
ERROR! Unexpected Exception: [Errno 104] Connection reset by peer
```

## Features

* Implements `ansible-galaxy install` command only.
* Supports looking up the latest role version from Ansible Galaxy.
* Efficient HTTP implementation, reuses HTTP connections.
* Fast, roles are installed in parallel making this implementation much faster.

    * For example a `requirements.yaml` with 30 roles took **2:32 mins** to
      install with `ansible-galaxy` and only **2.75 secs!** with `alt-galaxy`.

* Only supports installing roles from a role file.

    * Directly installing individual tar files is not supported.

* Only supports installing `tar.gz` packages.

    * Use of `git` or `hg` is not supported.

* Standalone, the application is written in [Go](https://golang.org/) and
  statically linked, so it has no dependencies other than the operating system.

## Usage guide

```bash
alt-galaxy install --role-file=/vagrant/provisioning/requirements.yml --roles-path=/etc/ansible/roles
```

## License

This software is licensed under the terms in the file named "[LICENSE](LICENSE)"
in the root directory of this project.

## Author Information

John Freeman

GantSign Ltd.
Company No. 06109112 (registered in England)
