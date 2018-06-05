from nozzlebeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Nozzlebeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        nozzlebeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("nozzlebeat is running"))
        exit_code = nozzlebeat_proc.kill_and_wait()
        assert exit_code == 0
