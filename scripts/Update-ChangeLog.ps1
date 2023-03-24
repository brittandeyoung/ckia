param (
    [Parameter(Mandatory=$true)]
    [System.Management.Automation.SemanticVersion]
    $ReleaseVersion
)

$LinkPattern   = @{
    FirstRelease  = "https://github.com/brittandeyoung/ckia/tree/v{CUR}"
    NormalRelease = "https://github.com/brittandeyoung/ckia/compare/v{PREV}..v{CUR}"
    Unreleased    = "https://github.com/brittandeyoung/ckia/compare/v{CUR}..HEAD"
}
Update-Changelog -ReleaseVersion $ReleaseVersion -LinkMode Automatic -LinkPattern $LinkPattern