"""MkDocs hooks: preview banner and safer external link attributes."""

from __future__ import annotations

import os
import re


def on_config(config):
    channel = os.environ.get("DOCS_CHANNEL", "stable")
    config.extra["preview_banner"] = channel == "preview"
    return config


def on_page_markdown(markdown, page, config, files):
    if config.extra.get("preview_banner"):
        banner = (
            "> **Preview documentation** — this branch is not release-ready. "
            "See the [stable docs](https://synthread.gitlab.io/proj/Klippventurer/) "
            "for the current manual-install line.\n\n"
        )
        markdown = banner + markdown
    return markdown


_EXTERNAL_LINK_RE = re.compile(
    r'<a href="(https?://[^"]+)"(?![^>]*\btarget=)',
)


def on_post_page(output, page, config):
    """Add rel and target on external links (issue #2)."""

    def _replace(match: re.Match[str]) -> str:
        href = match.group(1)
        return (
            f'<a href="{href}" target="_blank" rel="noopener noreferrer"'
        )

    return _EXTERNAL_LINK_RE.sub(_replace, output)
