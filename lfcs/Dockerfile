# LFCS Exam Practice Environment
# Based on Ubuntu 22.04 LTS

FROM ubuntu:22.04

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Remove man page excludes (Ubuntu minimized image excludes man pages by default)
RUN rm -f /etc/dpkg/dpkg.cfg.d/excludes* && \
    sed -i 's/path-exclude=\/usr\/share\/man\/\*/# path-exclude=\/usr\/share\/man\/*/g' /etc/dpkg/dpkg.cfg.d/* 2>/dev/null || true

# Install essential packages for LFCS exam preparation
# Based on official LFCS exam requirements and common Linux system administration tools
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    # Basic system utilities
    sudo \
    vim \
    nano \
    less \
    # Man pages (allowed resource during exam)
    man-db \
    manpages \
    manpages-dev \
    # Networking tools
    net-tools \
    iputils-ping \
    iproute2 \
    iptables \
    ufw \
    dnsutils \
    # Package management (apt, dpkg included by default in Ubuntu)
    apt-utils \
    # Compression tools
    tar \
    gzip \
    bzip2 \
    xz-utils \
    # Process and system monitoring
    htop \
    lsof \
    strace \
    procps \
    sysstat \
    # File utilities
    tree \
    acl \
    attr \
    # Service management
    systemd \
    systemd-sysv \
    cron \
    # Storage management
    lvm2 \
    parted \
    xfsprogs \
    # Services for practice
    nginx \
    postgresql-client \
    # Miscellaneous
    ca-certificates \
    curl \
    wget \
    git \
    bc \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Restore man pages for essential packages
RUN apt-get update && \
    apt-get install -y --reinstall \
    man-db \
    manpages \
    manpages-dev \
    sudo \
    systemd \
    cron \
    nginx \
    && mandb -q \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create users and groups for practice
RUN groupadd devteam && \
    groupadd operators && \
    useradd -m -s /bin/bash student && \
    useradd -m -s /bin/bash backup-user && \
    echo 'student:student' | chpasswd && \
    echo 'backup-user:backup' | chpasswd && \
    echo 'root:root' | chpasswd

# Setup sudo for student user
# LFCS exam uses: sudo -i to get root access
RUN echo 'student ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Create practice directories
RUN mkdir -p /opt/course && \
    mkdir -p /opt/mock/{1..10} && \
    mkdir -p /opt/backup && \
    mkdir -p /opt/important-data && \
    mkdir -p /backup && \
    mkdir -p /var/log/app && \
    mkdir -p /var/www/html && \
    chown student:student /opt/course && \
    chown student:student /opt/mock && \
    chown backup-user:backup-user /opt/backup && \
    chown www-data:www-data /var/www/html

# Create sample files for practice
RUN echo "Sample data 1" > /opt/important-data/file1.txt && \
    echo "Sample data 2" > /opt/important-data/file2.txt && \
    echo "127.0.0.1 - - [01/Jan/2024:12:00:00 +0000] \"GET /index.html HTTP/1.1\" 200 1234" >> /var/log/app/access.log && \
    echo "127.0.0.1 - - [01/Jan/2024:12:01:00 +0000] \"GET /missing.html HTTP/1.1\" 404 567" >> /var/log/app/access.log && \
    echo "127.0.0.1 - - [01/Jan/2024:12:02:00 +0000] \"GET /page.html HTTP/1.1\" 200 890" >> /var/log/app/access.log && \
    echo "127.0.0.1 - - [01/Jan/2024:12:03:00 +0000] \"GET /error.html HTTP/1.1\" 404 123" >> /var/log/app/access.log

# Create backup script
RUN echo '#!/bin/bash\ndate >> /tmp/backup.log' > /opt/backup/daily.sh && \
    chmod +x /opt/backup/daily.sh

# Set working directory
WORKDIR /home/student

# Switch to student user by default
USER student

# Start bash shell
CMD ["/bin/bash"]
