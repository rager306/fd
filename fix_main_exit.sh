cat << 'INNER_EOF' > patch_main_exit.diff
<<<<<<< SEARCH
		logger.Error("shutdown failed", "error", err)
		closeResource("redis", redisCache, logger)
		closeResource("local cache", localCache, logger)
		os.Exit(1)
	}
=======
		logger.Error("shutdown failed", "error", err)
		closeResource("redis", redisCache, logger)
		closeResource("local cache", localCache, logger)
		recoveryCancel()
		//nolint:gocritic // exitAfterDefer: cleanup happens before exit
		os.Exit(1)
	}
>>>>>>> REPLACE
INNER_EOF
