# typo3-staticfilecache-cleaner

A tool to clean static file caches for TYPO3.

## Why

If you're using the TYPO3 extension [staticfilecache](https://github.com/lochmueller/staticfilecache) with Nginx, you
can't rely on the generated .htaccess file to lead requests back to TYPO3 (when the cache has expired).
This tool removes them after they've been expired.

## Usage

```bash
$ typo3-staticfilecache-cleaner /var/www/html/typo3temp/tx_staticfilecache --dry-run
```

### With Docker

```bash
$ docker run --rm -v /var/www/html/typo3temp/tx_staticfilecache:/data atomicptr/typo3-staticfilecache-cleaner
```

## License

MIT