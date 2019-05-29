# Wordpress CLI Cloud Native Buildpack

If you have deployed Wordpress to Cloud Foundry then it maybe desirable to have the [`wp` CLI](https://wp-cli.org/) available when you `cf ssh` into your application container.

For example,

```plain
# cf ssh wordpress
$ /tmp/lifecycle/shell
$ wp plugin list
+-------------------+----------+--------+---------+
| name              | status   | update | version |
+-------------------+----------+--------+---------+
| aryo-activity-log | active   | none   | 2.5.2   |
| akismet           | inactive | none   | 4.1.2   |
| hello             | inactive | none   | 1.7.2   |
| s3-uploads        | active   | none   | 2.0.0   |
+-------------------+----------+--------+---------+
```

Learn more about the `wp` CLI in the [handbook](https://make.wordpress.org/cli/handbook/).

Inside `cf ssh`, the `wp` helper is already configured with the `--path=$HOME/htdocs` where Wordpress will be installed by the [php-buildpack](https://github.com/cloudfoundry/php-buildpack).