import gradle.VersionWatcher

final watcher = new VersionWatcher()

watcher.performIfNewRcReleaseIsAvailable {
    memcache.clearCacheForUri('/rc/feed')
}
