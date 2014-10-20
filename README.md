# php-classname-fixer

`php-classname-fixer` is a tool to help rewrite PHP class names.

It can be used to convert code to use [PHP namespaces](http://php.net/manual/en/language.namespaces.php) or more generally to rename classes.

You can install it with `go install github.com/mtibben/php-classname-fixer`.

## How it works

`php-classname-fixer` searches a directory for php files. It builds a map of classnames, then determines the implied PSR-0 classname from the location of the files. It then rewrites the classname and other PHP code to use the implied PSR-0 class name.

So for example, if I have a PHP class file called `MyClass.php`
```php
class OldVendor_OldNamespace_MyClass {
    public function getDate() {
        return new DateTime("2014-10-20");
    }
}
```

and I put that class in `src/NewVendor/NewNamespace/MyClass.php` and then run `php-classname-fixer src`, the file contents now look like this:

```php
namespace NewVendor\NewNamespace;

class MyClass {
    public function getDate() {
        return new \DateTime("2014-10-20");
    }
}
```

This is very useful if you need to move around large numbers of files around in your codebase - simply move the files to their desired location, and run `php-classname-fixer`.

## Notes

Regexes are used to rewrite code. To keep the regexes simple it assumed that all code being analysed is using the [PSR-2](http://www.php-fig.org/psr/psr-2/) code style. It is recommended that your code has been formatted that way before using this tool. ([PHP-CS-Fixer](https://github.com/fabpot/PHP-CS-Fixer) is an excellent tool to convert code to PSR-2)
